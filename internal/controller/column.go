package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/shared/form"
	"goat-cg/internal/service"
)


type ColumnController struct {
	columnService  service.ColumnService
	tableService service.TableService
	urlCheckService service.UrlCheckService
}


func NewColumnController() *ColumnController {
	columnService  := service.NewColumnService()
	tableService := service.NewTableService()
	urlCheckService := service.NewUrlCheckService()
	return &ColumnController{columnService , tableService, urlCheckService}
}


//GET /:project_cd/tables/:table_id/columns
func (ctr *ColumnController) columnsPage(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)

	table, _ := ctr.tableService.GetTable(tableId)
	columns, _ := ctr.columnService .GetColumns(tableId)

	c.HTML(200, "columns.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"table": table,
		"columns": columns,
	})
}


//GET /:project_cd/tables/:table_id/columns/new
func (ctr *ColumnController) createColumnPage(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)

	table, _ := ctr.tableService.GetTable(tableId)

	c.HTML(200, "column.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"table": table,
	})
}


//POST /:project_cd/tables/:table_id/columns/new
func (ctr *ColumnController) createColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)

	var form form.PostColumnsForm
	c.Bind(&form)
	result := ctr.columnService .CreateColumn(form.ToServInCreateColumn(tableId, userId))

	if result == service.CREATE_COLUMN_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/tables/%d/columns", c.Param("project_cd"), tableId))
		return
	}

	table, _ := ctr.tableService.GetTable(tableId)

	if result == service.CREATE_COLUMN_CONFLICT_INT {
		c.HTML(409, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "同一ColumnNameが既に登録されています",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "登録に失敗しました",
		})
	}
}


//GET /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) updateColumnPage(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlCheckService.CheckColumnIdAndGetColumnId(c, tableId)

	table, _ := ctr.tableService.GetTable(tableId)
	column, _ := ctr.columnService .GetColumn(columnId)

	c.HTML(200, "column.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"table": table,
		"column": column,
	})
}


//POST /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) updateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlCheckService.CheckColumnIdAndGetColumnId(c, tableId)

	var form form.PostColumnsForm
	c.Bind(&form)
	result := ctr.columnService .UpdateColumn(
		columnId, form.ToServInCreateColumn(tableId, userId),
	)

	if result == service.UPDATE_COLUMN_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/tables/%d/columns", c.Param("project_cd"), tableId))
		return
	}

	table, _ := ctr.tableService.GetTable(tableId)

	if result == service.UPDATE_COLUMN_CONFLICT_INT {
		c.HTML(409, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "同一ColumnNameが既に登録されています",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"commons": constant.Commons,
			"project_cd" : c.Param("project_cd"),
			"table": table,
			"column": form,
			"error": "更新に失敗しました",
		})
	}
}


//DELETE /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) deleteColumn(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlCheckService.CheckColumnIdAndGetColumnId(c, tableId)

	ctr.columnService .DeleteColumn(columnId)

	c.Redirect(303, fmt.Sprintf("/%s/tables/%d/columns", c.Param("project_cd"), tableId))

}


//GET /:project_cd/tables/:table_id/columns/:column_id/log
func (ctr *ColumnController) columnLogPage(c *gin.Context) {
	projectId := ctr.urlCheckService.CheckProjectCdAndGetProjectId(c)
	tableId := ctr.urlCheckService.CheckTableIdAndGetTableId(c, projectId)
	columnId := ctr.urlCheckService.CheckColumnIdAndGetColumnId(c, tableId)

	columnLog, _ := ctr.columnService .GetColumnLog(columnId)

	c.HTML(200, "columnlog.html", gin.H{
		"commons": constant.Commons,
		"project_cd" : c.Param("project_cd"),
		"columnlog": columnLog,
	})
}