package controller

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type tableController struct {
	tServ service.TableService
	urlServ service.UrlCheckService
}


func newTableController() *tableController {
	tServ := service.NewTableService()
	urlServ := service.NewUrlCheckService()
	return &tableController{tServ, urlServ}
}


//GET /:project_cd/tables
func (ctr *tableController) tablesPage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tables, _ := ctr.tServ.GetTables(projectId)

	c.HTML(200, "tables.html", gin.H{
		"commons": constant.Commons,
		"project_cd": c.Param("project_cd"), 
		"tables": tables,
	})
}


//GET /:project_cd/tables/new
func (ctr *tableController) createTablePage(c *gin.Context) {
	ctr.urlServ.CheckProjectCdAndGetProjectId(c)

	c.HTML(200, "table.html", gin.H{
		"commons": constant.Commons,
		"project_cd": c.Param("project_cd"), 
	})
}


//POST /:project_cd/tables/new
func (ctr *tableController) createTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")

	result := ctr.tServ.CreateTable(projectId, userId, tableName, tableNameLogical)

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
func (ctr *tableController) updateTablePage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

	table, _ := ctr.tServ.GetTable(tableId)

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
func (ctr *tableController) updateTable(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

	tableName := c.PostForm("table_name")
	tableNameLogical := c.PostForm("table_name_logical")
	delFlg, err := strconv.Atoi(c.PostForm("del_flg"))
	if err != nil || delFlg != 1 {
		delFlg = 0
	}

	result := ctr.tServ.UpdateTable(
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
func (ctr *tableController) deleteTable(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

	ctr.tServ.DeleteTable(tableId)

	c.Redirect(303, fmt.Sprintf("/%s/tables", c.Param("project_cd")))

}


//GET /:project_cd/tables/:table_id/log
func (ctr *tableController) tableLogPage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

	tableLog, _ := ctr.tServ.GetTableLog(tableId)

	c.HTML(200, "tablelog.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"tablelog": tableLog,
	})
}