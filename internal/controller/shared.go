package controller

import (
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/service"
)


/* 
URLパラメータ[:project_cd]がログインユーザのアクセス可能なプロジェクトかをチェック
アクセス不可 -> Redirect
アクセス可能 -> ProjectCd を ProjectIdに変換
*/
func CheckProjectCdAndGetProjectId(c *gin.Context) int {
    userId := jwt.GetUserId(c)
    projectCd := c.Param("project_cd")
    projectId := service.GetProjectId(userId, projectCd)

    if projectId == service.GET_PROJECT_ID_NOT_FOUND_INT {
        c.Redirect(303, "/projects")
        return -1
    }

    return projectId
}