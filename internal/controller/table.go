package controller

import (
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/shared/constant"
    "goat-cg/internal/service"
)


type tableController struct {
    tServ service.TableService
}


func newTableController() *tableController {
    tServ := service.NewTableService()
    return &tableController{tServ}
}


//GET /:project_cd/tables
func (ctr *tableController) tablesPage(c *gin.Context) {
    projectId := CheckProjectCdAndGetProjectId(c)
    tables, _ := ctr.tServ.GetTables(projectId)

    c.HTML(200, "tables.html", gin.H{
        "commons": constant.Commons,
        "tables": tables,
    })
}


//POST /:project_cd/tables
func (ctr *tableController) postTable(c *gin.Context) {
    userId := jwt.GetUserId(c)
    projectId := CheckProjectCdAndGetProjectId(c)
    tn := c.PostForm("table_name")
    tnl := c.PostForm("table_name_logical")

    ctr.tServ.CreateTable(userId, projectId, tn, tnl)

    tables, _ := ctr.tServ.GetTables(projectId)

    c.HTML(200, "tables.html", gin.H{
        "commons": constant.Commons,
        "tables": tables,
    })  
}