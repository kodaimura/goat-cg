package controller

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/service"
)


/* 
URLパラメータ/:project_cd がログインユーザのアクセス可能なProjectかをチェック
アクセス不可 -> Redirect
アクセス可能 -> ProjectCd を ProjectIdに変換
*/
func CheckProjectCdAndGetProjectId(c *gin.Context) int {
    userId := jwt.GetUserId(c)
    projectCd := c.Param("project_cd")

    pServ := service.NewProjectService()
    projectId := pServ.GetProjectId(userId, projectCd)

    if projectId == service.GET_PROJECT_ID_NOT_FOUND_INT {
        c.Redirect(303, "/projects")
        return -1
    }

    return projectId
}


/* 
URLパラメータ /:tableId がProject管理のTableかをチェック
アクセス不可 -> Redirect
アクセス可能 -> TableId を抽出
*/
func CheckTableIdAndGetTableId(c *gin.Context, projectId int) int {
    tableId, err := strconv.Atoi(c.Param("table_id"))
    if err != nil {
        c.Redirect(303, "/projects")
        return -1
    }

    tServ := service.NewTableService()
    _, err = tServ.GetTable(projectId, tableId)

    if err != nil {
        c.Redirect(303, "/projects")
        return -1
    }
    return tableId
}