package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/errs"
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
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	columns, _ := ctr.columnService.GetColumns(table.TableId)

	c.HTML(200, "columns.html", gin.H{
		"project" : project,
		"table": table,
		"columns": columns,
	})
}


//GET /:username/:project_name/tables/:table_id/columns/new
func (ctr *ColumnController) CreateColumnPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	c.HTML(200, "column.html", gin.H{
		"project" : project,
		"table": table,
	})
}


//POST /:username/:project_name/tables/:table_id/columns/new
func (ctr *ColumnController) CreateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	var form form.PostColumnsForm
	c.Bind(&form)
	err := ctr.columnService.CreateColumn(form.ToServInCreateColumn(table.TableId, userId))

	if err == nil {
		c.Redirect(303, fmt.Sprintf("/%s/%s/tables/%s/columns", c.Param("username"), c.Param("project_name"), c.Param("table_id")))
		return
	}

	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "column.html", gin.H{
			"project" : project,
			"table": table,
			"column": form,
			"error": "同一ColumnNameが既に登録されています",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"project" : project,
			"table": table,
			"column": form,
			"error": "登録に失敗しました",
		})
	}
}


//GET /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) UpdateColumnPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)
	column := c.Keys["column"].(model.Column)

	c.HTML(200, "column.html", gin.H{
		"project" : project,
		"table": table,
		"column": column,
	})
}


//POST /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) UpdateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)
	column := c.Keys["column"].(model.Column)

	var form form.PostColumnsForm
	c.Bind(&form)
	form.ColumnId = column.ColumnId
	err := ctr.columnService.UpdateColumn(form.ToServInCreateColumn(table.TableId, userId))

	if err == nil {
		c.Redirect(303, fmt.Sprintf("/%s/%s/tables/%s/columns", c.Param("username"), c.Param("project_name"), c.Param("table_id")))
		return
	}

	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "column.html", gin.H{
			"project" : project,
			"table": table,
			"column": form,
			"error": "同一ColumnNameが既に登録されています",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"project" : project,
			"table": table,
			"column": form,
			"error": "更新に失敗しました",
		})
	}
}


//DELETE /:project_cd/tables/:table_id/columns/:column_id
func (ctr *ColumnController) DeleteColumn(c *gin.Context) {
	column := c.Keys["column"].(model.Column)

	ctr.columnService.DeleteColumn(column.ColumnId)

	c.Redirect(303, fmt.Sprintf("/%s/%s/tables/%s/columns", c.Param("username"), c.Param("project_name"), c.Param("table_id")))

}


//GET /:project_cd/tables/:table_id/columns/:column_id/log
func (ctr *ColumnController) ColumnLogPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	column := c.Keys["column"].(model.Column)
	columnLog, _ := ctr.columnService.GetColumnLog(column.ColumnId)

	c.HTML(200, "columnlog.html", gin.H{
		"project" : project,
		"columnlog": columnLog,
	})
}