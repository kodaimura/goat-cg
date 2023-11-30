package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/shared/form"
	"goat-cg/internal/service"
	"goat-cg/internal/model"
)


type ColumnController struct {
	columnService  service.ColumnService
	tableService service.TableService
}


func NewColumnController() *ColumnController {
	columnService  := service.NewColumnService()
	tableService := service.NewTableService()
	return &ColumnController{columnService , tableService}
}


//GET /:username/:project_name/tables/:table_id/columns
func (ctr *ColumnController) ColumnsPage(c *gin.Context) {
	t := c.Keys["table"].(model.Table)

	columns, _ := ctr.columnService.GetColumns(t.TableId)

	c.HTML(200, "columns.html", gin.H{
		"commons": constant.Commons,
		"project_name" : c.Param("project_name"),
		"table": t,
		"columns": columns,
	})
}


//GET /:username/:project_name/tables/:table_id/columns/new
func (ctr *ColumnController) CreateColumnPage(c *gin.Context) {
	t := c.Keys["table"].(model.Table)

	c.HTML(200, "column.html", gin.H{
		"commons": constant.Commons,
		"project_name" : c.Param("project_name"),
		"table": t,
	})
}


//POST /:username/:project_name/tables/:table_id/columns/new
func (ctr *ColumnController) CreateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	t := c.Keys["table"].(model.Table)

	var form form.PostColumnsForm
	c.Bind(&form)
	result := ctr.columnService.CreateColumn(form.ToServInCreateColumn(t.TableId, userId))

	if result == service.CREATE_COLUMN_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/%s/tables/%d/columns", c.Param("username"), c.Param("project_name"), t.TableId))
		return
	}

	if result == service.CREATE_COLUMN_CONFLICT_INT {
		c.HTML(409, "column.html", gin.H{
			"commons": constant.Commons,
			"project_name" : c.Param("project_name"),
			"table": t,
			"column": form,
			"error": "同一ColumnNameが既に登録されています",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"commons": constant.Commons,
			"project_name" : c.Param("project_name"),
			"table": t,
			"column": form,
			"error": "登録に失敗しました",
		})
	}
}


//GET /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) UpdateColumnPage(c *gin.Context) {
	userId := jwt.GetUserId(c)
	t := c.Keys["table"].(model.Table)
	col := c.Keys["column"].(model.Column)

	c.HTML(200, "column.html", gin.H{
		"commons": constant.Commons,
		"project_name" : c.Param("project_name"),
		"table": t,
		"column": col,
	})
}


//POST /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) UpdateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	t := c.Keys["table"].(model.Table)

	var form form.PostColumnsForm
	c.Bind(&form)
	form.ColumnId = c.Param("column_id")
	result := ctr.columnService.UpdateColumn(
		form.ToServInCreateColumn(t.TableId, userId),
	)

	if result == service.UPDATE_COLUMN_SUCCESS_INT {
		c.Redirect(303, fmt.Sprintf("/%s/%s/tables/%d/columns", c.Param("username"), c.Param("project_name"), t.TableId))
		return
	}

	if result == service.UPDATE_COLUMN_CONFLICT_INT {
		c.HTML(409, "column.html", gin.H{
			"commons": constant.Commons,
			"project_name" : c.Param("project_name"),
			"table": t,
			"column": form,
			"error": "同一ColumnNameが既に登録されています",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"commons": constant.Commons,
			"project_name" : c.Param("project_name"),
			"table": t,
			"column": form,
			"error": "更新に失敗しました",
		})
	}
}


//DELETE /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) DeleteColumn(c *gin.Context) {
	t := c.Keys["table"].(model.Table)

	ctr.columnService.DeleteColumn(c.Param("column_id"))

	c.Redirect(303, fmt.Sprintf("/%s/%s/tables/%d/columns", c.Param("username"), c.Param("project_name"), t.TableId))

}


//GET /:project_cd/tables/:table_id/columns/:column_id/log
func (ctr *ColumnController) ColumnLogPage(c *gin.Context) {
	columnLog, _ := ctr.columnService.GetColumnLog(c.Param("column_id"))

	c.HTML(200, "columnlog.html", gin.H{
		"commons": constant.Commons,
		"project_name" : c.Param("project_name"),
		"columnlog": columnLog,
	})
}