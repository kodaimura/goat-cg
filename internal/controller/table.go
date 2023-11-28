package controller

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type TableController struct {
	tableService service.TableService
	urlCheckService service.UrlCheckService
}


func NewTableController() *TableController {
	tableService := service.NewTableService()
	urlCheckService := service.NewUrlCheckService()
	return &TableController{tableService, urlCheckService}
}


//GET /:project_cd/tables
func (ctr *TableController) tablesPage(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tables, _ := ctr.tableService.GetTables(projectId)

	c.HTML(200, "tables.html", gin.H{
		"commons": constant.Commons,
		"project_cd": c.Param("project_cd"), 
		"tables": tables,
	})
}


//GET /:project_cd/tables/new
func (ctr *TableController) createTablePage(c *gin.Context) {
	ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)

	c.HTML(200, "table.html", gin.H{
		"commons": constant.Commons,
		"project_cd": c.Param("project_cd"), 
	})
}


//POST /:project_cd/tables/new
func (ctr *TableController) createTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")

	result := ctr.tableService.CreateTable(projectId, userId, tableName, tableNameLogical)

	if result == service.CREATE_TABLE_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/tables", c.Param("project_cd")))
		return
	}

	if result == service.CREATE_TABLE_CONFLICT_INT {
		c.HTML(409, "table.html", gin.H{
			"commons": constant.Commons,
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"error": "同一TableNameが既に登録されています",
		})
	} else {
		c.HTML(500, "table.html", gin.H{
			"commons": constant.Commons,
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"error": "登録に失敗しました",
		})
	}
}


//GET /:project_cd/tables/:table_id
func (ctr *TableController) updateTablePage(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)

	table, _ := ctr.tableService.GetTable(tableId)

	c.HTML(200, "table.html", gin.H{
		"commons": constant.Commons,
		"project_cd": c.Param("project_cd"),
		"table_id": tableId,
		"table_name": table.TableName,
		"table_name_logical": table.TableNameLogical,
		"del_flg": table.DelFlg,
	})
}


//POST /:project_cd/tables/:table_id
func (ctr *TableController) updateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")
	delFlg, err := strconv.Atoi(c.PostForm("del_flg"))
	if err != nil || delFlg != 1 {
		delFlg = 0
	}

	result := ctr.tableService.UpdateTable(
		projectId, tableId, userId, tableName, tableNameLogical, delFlg,
	)

	if result == service.UPDATE_TABLE_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/tables", c.Param("project_cd")))

	} else if result == service.UPDATE_TABLE_CONFLICT_INT {
		c.HTML(409, "table.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table_id": tableId,
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"del_flg": delFlg,
			"error": "同一TableNameが既に登録されています",
		})
	} else {
		c.HTML(500, "table.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table_id": tableId,
			"table_name": tableName,
			"table_name_logical": tableNameLogical,
			"del_flg": delFlg,
			"error": "更新に失敗しました",
		})
	}
}


//DELETE /:project_cd/tables/:table_id
func (ctr *TableController) deleteTable(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)

	ctr.tableService.DeleteTable(tableId)

	c.Redirect(303, fmt.Sprintf("/%s/tables", c.Param("project_cd")))

}


//GET /:project_cd/tables/:table_id/log
func (ctr *TableController) tableLogPage(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)

	tableLog, _ := ctr.tableService.GetTableLog(tableId)

	c.HTML(200, "tablelog.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"tablelog": tableLog,
	})
}