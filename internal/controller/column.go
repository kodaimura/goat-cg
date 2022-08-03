package controller

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "goat-cg/internal/constant"
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
    tableId, err := strconv.Atoi(c.Param("table_id"))

    if err != nil {
        c.Redirect(303, "/projects")
        return
    }

    table, err := ctr.tServ.GetTable(projectId, tableId)

    if err != nil {
        c.Redirect(303, "/projects")
        return
    }

    columns, _ := ctr.cServ.GetColumns(tableId)

    c.HTML(200, "columns.html", gin.H{
        "commons": constant.Commons,
        "table": table,
        "columns": columns,
    })
}


//POST /:project_cd/tables/:table_id/columns
func (ctr *columnController) postColumns(c *gin.Context) {
    c.HTML(200, "tables.html", gin.H{
        "commons": constant.Commons,
    })  
}