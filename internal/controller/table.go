package controller

import (
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/constant"
    "goat-cg/internal/service"
    //"goat-cg/internal/model/queryservice"
)


func setTableRoute(r *gin.Engine) {
    tc := newTableController()

    p := r.Group("/:project_cd", jwt.JwtAuthMiddleware())
    {
        p.GET("/tables", tc.tablesPage)
        p.POST("/tables", tc.postTable)
    }
}

type tableController struct {
    tServ service.TableService
}


func newTableController() *tableController {
    tServ := service.NewTableService()
    return &tableController{tServ}
}


//GET /:project_cd/tables
func (ctr *tableController) tablesPage(c *gin.Context) {
    projectId := checkProjectCdAndConvertToProjectId(c)
    tables, _ := ctr.tServ.GetTables(projectId)

    c.HTML(200, "tables.html", gin.H{
        "commons": constant.Commons,
        "tables": tables,
    })
}


//POST /:project_cd/tables
func (ctr *tableController) postTable(c *gin.Context) {
    userId := jwt.GetUserId(c)
    projectId := checkProjectCdAndConvertToProjectId(c)
    tn := c.PostForm("table_name")
    tnl := c.PostForm("table_name_logical")

    ctr.tServ.CreateTable(userId, projectId, tn, tnl)

    tables, _ := ctr.tServ.GetTables(projectId)

    c.HTML(200, "tables.html", gin.H{
        "commons": constant.Commons,
        "tables": tables,
    })  
}

/* 
URLパラメータ[:project_cd]がログインユーザのアクセス可能なプロジェクトかをチェック
アクセス不可 -> Redirect
アクセス可能 -> ProjectCd を ProjectIdに変換
*/
func checkProjectCdAndConvertToProjectId(c *gin.Context) int {
    userId := jwt.GetUserId(c)
    projectCd := c.Param("project_cd")
    projectId := service.GetProjectId(userId, projectCd)

    if projectId == service.GET_PROJECT_ID_NOT_FOUND_INT {
        c.Redirect(303, "/projects")
        return -1
    }

    return projectId
}