package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/shared/form"
	"goat-cg/internal/service"
)


type columnController struct {
	cServ service.ColumnService
	tServ service.TableService
	urlServ service.UrlCheckService
}


func newColumnController() *columnController {
	cServ := service.NewColumnService()
	tServ := service.NewTableService()
	urlServ := service.NewUrlCheckService()
	return &columnController{cServ, tServ, urlServ}
}


//GET /:project_cd/tables/:table_id/columns
func (ctr *columnController) columnsPage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

	table, _ := ctr.tServ.GetTable(tableId)
	columns, _ := ctr.cServ.GetColumns(tableId)

	c.HTML(200, "columns.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"table": table,
		"columns": columns,
	})
}


//GET /:project_cd/tables/:table_id/columns/new
func (ctr *columnController) createColumnPage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

	table, _ := ctr.tServ.GetTable(tableId)

	c.HTML(200, "column.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"table": table,
	})
}


//POST /:project_cd/tables/:table_id/columns/new
func (ctr *columnController) createColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

	var form form.PostColumnsForm
	c.Bind(&form)
	result := ctr.cServ.CreateColumn(form.ToServInCreateColumn(tableId, userId))

	if result == service.CREATE_COLUMN_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/tables/%d/columns", c.Param("project_cd"), tableId))
		return
	}

	table, _ := ctr.tServ.GetTable(tableId)

	if result == service.CREATE_COLUMN_CONFLICT_INT {
		c.HTML(409, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "??????ColumnName?????????????????????????????????",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "???????????????????????????",
		})
	}
}


//GET /:project_cd/tables/:table_id/columns/:column_id
func (ctr *columnController) updateColumnPage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlServ.CheckColumnIdAndGetColumnId(c, tableId)

	table, _ := ctr.tServ.GetTable(tableId)
	column, _ := ctr.cServ.GetColumn(columnId)

	c.HTML(200, "column.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"table": table,
		"column": column,
	})
}


//POST /:project_cd/tables/:table_id/columns/:column_id
func (ctr *columnController) updateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlServ.CheckColumnIdAndGetColumnId(c, tableId)

	var form form.PostColumnsForm
	c.Bind(&form)
	result := ctr.cServ.UpdateColumn(
		columnId, form.ToServInCreateColumn(tableId, userId),
	)

	if result == service.UPDATE_COLUMN_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/tables/%d/columns", c.Param("project_cd"), tableId))
		return
	}

	table, _ := ctr.tServ.GetTable(tableId)

	if result == service.UPDATE_COLUMN_CONFLICT_INT {
		c.HTML(409, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "??????ColumnName?????????????????????????????????",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "???????????????????????????",
		})
	}
}


//DELETE /:project_cd/tables/:table_id/columns/:column_id
func (ctr *columnController) deleteColumn(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlServ.CheckColumnIdAndGetColumnId(c, tableId)

	ctr.cServ.DeleteColumn(columnId)

	c.Redirect(303, fmt.Sprintf("/%s/tables/%d/columns", c.Param("project_cd"), tableId))

}


//GET /:project_cd/tables/:table_id/columns/:column_id/log
func (ctr *columnController) columnLogPage(c *gin.Context) {
	projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlServ.CheckColumnIdAndGetColumnId(c, tableId)

	columnLog, _ := ctr.cServ.GetColumnLog(columnId)

	c.HTML(200, "columnlog.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"columnlog": columnLog,
	})
}