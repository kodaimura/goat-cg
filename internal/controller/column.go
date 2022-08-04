package controller

import (
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/shared/constant"
    "goat-cg/internal/shared/form"
    "goat-cg/internal/service"
)


type columnController struct {
    cServ service.ColumnService
    tServ service.TableService
}


func newColumnController() *columnController {
    cServ := service.NewColumnService()
    tServ := service.NewTableService()
    return &columnController{cServ, tServ}
}


//GET /:project_cd/tables/:table_id/columns
func (ctr *columnController) columnsPage(c *gin.Context) {
    projectId := CheckProjectCdAndGetProjectId(c)
    tableId := CheckTableIdAndGetTableId(c, projectId)

    tn := c.PostForm("table_name")
    tnl := c.PostForm("table_name_logical")

    ctr.tServ.CreateTable(jwt.GetUserId(c), projectId, tn, tnl)

    table, _ := ctr.tServ.GetTable(projectId, tableId)
    columns, _ := ctr.cServ.GetColumns(tableId)

    c.HTML(200, "columns.html", gin.H{
        "commons": constant.Commons,
        "table": table,
        "columns": columns,
    })
}


//POST /:project_cd/tables/:table_id/columns
func (ctr *columnController) postColumns(c *gin.Context) {
    projectId := CheckProjectCdAndGetProjectId(c)
    tableId := CheckTableIdAndGetTableId(c, projectId)

    var form *form.PostColumnsForm
    c.BindJSON(&form)

    ctr.cServ.CreateColumn(form.ToServInCreateColumn(tableId, jwt.GetUserId(c)))

    c.Redirect(303, c.Request.URL.Path)
}

