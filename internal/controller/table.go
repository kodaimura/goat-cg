package controller

import (
    "fmt"
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

    ctr.tServ.CreateTable(projectId, userId, tableName, tableNameLogical)

    c.Redirect(303, fmt.Sprintf("/%s/tables", c.Param("project_cd")))
}


//GET /:project_cd/tables/:table_id
func (ctr *tableController) updateTablePage(c *gin.Context) {
    projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
    tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

    table, _ := ctr.tServ.GetTable(tableId)

    c.HTML(200, "table.html", gin.H{
        "commons": constant.Commons,
        "project_cd": c.Param("project_cd"),
        "table": table,
    })
}


//POST /:project_cd/tables/:table_id
func (ctr *tableController) updateTable(c *gin.Context) {
    userId := jwt.GetUserId(c)
    projectId := ctr.urlServ.CheckProjectCdAndGetProjectId(c)
    tableId := ctr.urlServ.CheckTableIdAndGetTableId(c, projectId)

    tableName := c.PostForm("table_name")
    tableNameLogical := c.PostForm("table_name_logical")

    ctr.tServ.UpdateTable(tableId, userId, tableName, tableNameLogical)

    c.Redirect(303, fmt.Sprintf("/%s/tables", c.Param("project_cd")))
}