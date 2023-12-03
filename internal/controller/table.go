package controller

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/service"
	"goat-cg/internal/model"
)


type TableController struct {
	projectService service.ProjectService
	tableService service.TableService
}


func NewTableController() *TableController {
	projectService := service.NewProjectService()
	tableService := service.NewTableService()
	return &TableController{projectService, tableService}
}


//GET /:username/:project_name or /:username/:project_name/tables
func (ctr *TableController) TablesPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)

	tables, _ := ctr.tableService.GetTables(project.ProjectId)
	c.HTML(200, "tables.html", gin.H{
		"project": project, 
		"tables": tables,
	})
}


//GET /:username/:project_name/tables/new
func (ctr *TableController) CreateTablePage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)

	c.HTML(200, "table.html", gin.H{
		"project": project, 
	})
}


//POST /:username/:project_name/tables/new
func (ctr *TableController) CreateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")

	err := ctr.tableService.CreateTable(project.ProjectId, userId, tableName, tableNameLogical)

	if err == nil {
		c.Redirect(303, fmt.Sprintf("/%s/%s", c.Param("username"), c.Param("project_name")))
		return
	} 

	var table model.Table
	table.TableName = tableName
	table.TableNameLogical = tableNameLogical

	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "table.html", gin.H{
			"project": project, 
			"table": table,
			"error": "同一TableNameが既に登録されています",
		})

	} else {
		c.HTML(500, "table.html", gin.H{
			"project": project, 
			"table": table,
			"error": "登録に失敗しました",
		})
	}
}


//GET /:username/:project_name/tables/:table_id
func (ctr *TableController) UpdateTablePage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	c.HTML(200, "table.html", gin.H{
		"project": project, 
		"table": table,
	})
}


//POST /:username/:project_name/tables/:table_id
func (ctr *TableController) UpdateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")
	delFlg, err := strconv.Atoi(c.PostForm("del_flg"))

	if err != nil || delFlg != 1 {
		delFlg = 0
	}

	err = ctr.tableService.UpdateTable(
		project.ProjectId, table.TableId, userId, tableName, tableNameLogical, delFlg,
	)

	if err == nil {
		c.Redirect(303, fmt.Sprintf("/%s/%s", c.Param("username"), c.Param("project_name")))
		return
	} 

	table.TableName = tableName
	table.TableNameLogical = tableNameLogical
	table.DelFlg = delFlg

	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "table.html", gin.H{
			"project": project, 
			"taple": table,
			"error": "同一TableNameが既に登録されています",
		})
	} else {
		c.HTML(500, "table.html", gin.H{
			"project": project, 
			"taple": table,
			"error": "更新に失敗しました",
		})
	}
}


//DELETE /:username/:project_name/tables/:table_id
func (ctr *TableController) DeleteTable(c *gin.Context) {
	table := c.Keys["table"].(model.Table)
	ctr.tableService.DeleteTable(table.TableId)

	c.Redirect(303, fmt.Sprintf("/%s/%s/tables", c.Param("username"), c.Param("project_name")))

}


//GET /:username/:project_name/tables/:table_id/log
func (ctr *TableController) TableLogPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)
	tableLog, _ := ctr.tableService.GetTableLog(table.TableId)

	c.HTML(200, "tablelog.html", gin.H{
		"project": project, 
		"tablelog": tableLog,
	})
}