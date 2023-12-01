package controller

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
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


//GET /:username/:project_name/tables
func (ctr *TableController) TablesPage(c *gin.Context) {
	p := c.Keys["project_id"].(model.Project)

	tables, _ := ctr.tableService.GetTables(p.ProjectId)
	c.HTML(200, "tables.html", gin.H{
		"commons": constant.Commons,
		"project_name": c.Param("project_name"), 
		"tables": tables,
	})
}


//GET /:username/:project_name/tables/new
func (ctr *TableController) CreateTablePage(c *gin.Context) {

	c.HTML(200, "table.html", gin.H{
		"commons": constant.Commons,
		"project_name": c.Param("project_name"), 
	})
}


//POST /:username/:project_name/tables/new
func (ctr *TableController) CreateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectName := c.Param("project_name")
	p := c.Keys["project"].(model.Project)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")

	result := ctr.tableService.CreateTable(p.ProjectId, userId, tableName, tableNameLogical)

	if result == service.CREATE_TABLE_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/%s/tables", c.Param("username"), projectName))
		return
	}

	if result == service.CREATE_TABLE_CONFLICT_INT {
		c.HTML(409, "table.html", gin.H{
			"commons": constant.Commons,
			"project_name": projectName, 
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"error": "同一TableNameが既に登録されています",
		})
	} else {
		c.HTML(500, "table.html", gin.H{
			"commons": constant.Commons,
			"project_name": projectName, 
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"error": "登録に失敗しました",
		})
	}
}


//GET /:username/:project_name/tables/:table_id
func (ctr *TableController) UpdateTablePage(c *gin.Context) {
	t := c.Keys["table"].(model.Table)

	c.HTML(200, "table.html", gin.H{
		"commons": constant.Commons,
		"project_name": c.Param("project_name"), 
		"table_id": t.TableId,
		"table_name": t.TableName,
		"table_name_logical": t.TableNameLogical,
		"del_flg": t.DelFlg,
	})
}


//POST /:username/:project_name/tables/:table_id
func (ctr *TableController) UpdateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	p := c.Keys["project"].(model.Project)
	t := c.Keys["table"].(model.Table)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")
	delFlg, err := strconv.Atoi(c.PostForm("del_flg"))
	if err != nil || delFlg != 1 {
		delFlg = 0
	}

	result := ctr.tableService.UpdateTable(
		p.ProjectId, t.TableId, userId, tableName, tableNameLogical, delFlg,
	)

	if result == service.UPDATE_TABLE_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/%s/tables", c.Param("username"), p.ProjectName))

	} else if result == service.UPDATE_TABLE_CONFLICT_INT {
		c.HTML(409, "table.html", gin.H{
			"commons": constant.Commons,
			"project_name": c.Param("project_name"), 
			"table_id": t.TableId,
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"del_flg": delFlg,
			"error": "同一TableNameが既に登録されています",
		})
	} else {
		c.HTML(500, "table.html", gin.H{
			"commons": constant.Commons,
			"project_name": c.Param("project_name"), 
			"table_id": t.TableId,
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"del_flg": delFlg,
			"error": "更新に失敗しました",
		})
	}
}


//DELETE /:username/:project_name/tables/:table_id
func (ctr *TableController) DeleteTable(c *gin.Context) {
	p := c.Keys["project"].(model.Project)
	t := c.Keys["table"].(model.Table)
	ctr.tableService.DeleteTable(t.TableId)

	c.Redirect(303, fmt.Sprintf("/%s/%s/tables", c.Param("username"), p.ProjectName))

}


//GET /:username/:project_name/tables/:table_id/log
func (ctr *TableController) TableLogPage(c *gin.Context) {
	t := c.Keys["table"].(model.Table)
	tableLog, _ := ctr.tableService.GetTableLog(t.TableId)

	c.HTML(200, "tablelog.html", gin.H{
		"commons": constant.Commons,
		"project_name": c.Param("project_name"), 
		"tablelog": tableLog,
	})
}