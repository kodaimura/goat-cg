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
func (cc *ColumnController) ColumnsPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	columns, _ := cc.columnService.GetColumns(table.TableId)

	c.HTML(200, "columns.html", gin.H{
		"project": project,
		"table": table,
		"columns": columns,
	})
}


//GET /:username/:project_name/tables/:table_id/columns/new
func (cc *ColumnController) CreateColumnPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	c.HTML(200, "column.html", gin.H{
		"project": project,
		"table": table,
	})
}


//POST /:username/:project_name/tables/:table_id/columns/new
func (cc *ColumnController) CreateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)

	var form form.PostColumn
	c.Bind(&form)
	err := cc.columnService.CreateColumn(form.ToCreateColumn(table.TableId, userId))

	if err == nil {
		c.Redirect(303, fmt.Sprintf(
			"/%s/%s/tables/%s/columns", 
			c.Param("username"), c.Param("project_name"), c.Param("table_id"),
		))
		return
	}

	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "column.html", gin.H{
			"project": project,
			"table": table,
			"column": form,
			"error": "ColumnName must be unique.",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"project": project,
			"table": table,
			"column": form,
			"error": "error occurred.",
		})
	}
}


//GET /:project_cd/tables/:table_id/columns/:column_id
func (cc *ColumnController) UpdateColumnPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)
	column := c.Keys["column"].(model.Column)

	c.HTML(200, "column.html", gin.H{
		"project": project,
		"table": table,
		"column": column,
	})
}


//POST /:project_cd/tables/:table_id/columns/:column_id
func (cc *ColumnController) UpdateColumn(c *gin.Context) {
	userId := jwt.GetUserId(c)
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)
	column := c.Keys["column"].(model.Column)

	var form form.PostColumn
	c.Bind(&form)
	form.ColumnId = column.ColumnId
	err := cc.columnService.UpdateColumn(form.ToCreateColumn(table.TableId, userId))

	if err == nil {
		c.Redirect(303, fmt.Sprintf(
			"/%s/%s/tables/%s/columns", 
			c.Param("username"), c.Param("project_name"), c.Param("table_id"),
		))
		return
	}

	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "column.html", gin.H{
			"project": project,
			"table": table,
			"column": form,
			"error": "ColumnName must be unique.",
		})
	} else {
		c.HTML(500, "column.html", gin.H{
			"project": project,
			"table": table,
			"column": form,
			"error": "error occurred.",
		})
	}
}


//DELETE /:project_cd/tables/:table_id/columns/:column_id
func (cc *ColumnController) DeleteColumn(c *gin.Context) {
	column := c.Keys["column"].(model.Column)

	if cc.columnService.DeleteColumn(column.ColumnId) != nil {
		c.JSON(500, gin.H{})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}


//GET /:project_cd/tables/:table_id/columns/:column_id/log
func (cc *ColumnController) ColumnLogPage(c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	table := c.Keys["table"].(model.Table)
	column := c.Keys["column"].(model.Column)
	columnLog, _ := cc.columnService.GetColumnLog(column.ColumnId)

	c.HTML(200, "columnlog.html", gin.H{
		"project": project,
		"table": table,
		"columnlog": columnLog,
	})
}