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
func (tc *TableController) TablesPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)

	tables, _ := tc.tableService.GetTables(project.ProjectId)
	c.HTML(200, "tables.html", gin.H{
		"username": jwt.GetUsername(c),
		"project": project, 
		"tables": tables,
	})
}


//GET /:username/:project_name/tables/new
func (tc *TableController) CreateTablePage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)

	c.HTML(200, "table.html", gin.H{
		"project": project, 
	})
}


//POST /:username/:project_name/tables/new
func (tc *TableController) CreateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")

	err := tc.tableService.CreateTable(project.ProjectId, userId, tableName, tableNameLogical)

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
			"error": "TableName must be unique.",
		})

	} else {
		c.HTML(500, "table.html", gin.H{
			"project": project, 
			"table": table,
			"error": "error occurred.",
		})
	}
}


//GET /:username/:project_name/tables/:table_id
func (tc *TableController) UpdateTablePage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	c.HTML(200, "table.html", gin.H{
		"project": project, 
		"table": table,
	})
}


//POST /:username/:project_name/tables/:table_id
func (tc *TableController) UpdateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")
	delFlg, err := strconv.Atoi(c.PostForm("del_flg"))

	if err != nil || delFlg != 1 {
		delFlg = 0
	}

	err = tc.tableService.UpdateTable(
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
			"error": "TableName must be unique.",
		})
	} else {
		c.HTML(500, "table.html", gin.H{
			"project": project, 
			"taple": table,
			"error": "error occurred.",
		})
	}
}


//DELETE /:username/:project_name/tables/:table_id
func (tc *TableController) DeleteTable(c *gin.Context) {
	table := c.Keys["table"].(model.Table)

	if tc.tableService.DeleteTable(table.TableId) != nil {
		c.JSON(500, gin.H{"error": "error occurred."})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}


//GET /:username/:project_name/tables/:table_id/log
func (tc *TableController) TableLogPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)
	tableLog, _ := tc.tableService.GetTableLog(table.TableId)

	c.HTML(200, "tablelog.html", gin.H{
		"project": project, 
		"tablelog": tableLog,
	})
}