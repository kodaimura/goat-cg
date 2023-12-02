package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type ProjectController struct {
	projectService  service.ProjectService
}


func NewProjectController() *ProjectController {
	projectService  := service.NewProjectService()
	return &ProjectController{projectService}
}


//GET /:username or /:username/projects
func (ctr *ProjectController) ProjectsPage(c *gin.Context) {
	userId := jwt.GetUserId(c)
	username := jwt.GetUsername(c)

	if c.Param("username") != username {
		c.HTML(404, "404error.html", gin.H{})
		c.Abort()
		return
	}

	projects, _ := ctr.projectService.GetProjects(userId)
	member_projects, _ := ctr.projectService.GetMemberProjects(userId)

	c.HTML(200, "index.html", gin.H{
		"commons": constant.Commons,
		"username": username,
		"projects": projects,
		"member_projects": member_projects,
	})
}


//GET /:username/projects/new
func (ctr *ProjectController) CreateProjectPage(c *gin.Context) {
	username := jwt.GetUsername(c)

	if c.Param("username") != username {
		c.HTML(404, "404error.html", gin.H{})
		c.Abort()
		return
	}
	
	c.HTML(200, "project.html", gin.H{
		"commons": constant.Commons,
		"username": username,
	})
}


//POST /:username/projects
func (ctr *ProjectController) CreateProject(c *gin.Context) {
	userId := jwt.GetUserId(c)
	username := jwt.GetUsername(c)

	if c.Param("username") != username {
		c.HTML(404, "404error.html", gin.H{})
		c.Abort()
		return
	}

	projectName := c.PostForm("project_name")
	projectMemo := c.PostForm("project_memo")
	err := ctr.projectService.CreateProject(userId, username, projectName, projectMemo)
	
	if err == nil {
		c.Redirect(303, fmt.Sprintf("/%s", username))

	} else if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "project.html", gin.H{
			"commons": constant.Commons,
			"username": username,
			"error": "プロジェクト名が重複して使われています",
			"project_name": projectName,
			"project_memo": projectMemo,
		})

	} else {
		c.HTML(500, "project.html", gin.H{
			"commons": constant.Commons,
			"username": username,
			"error": "登録に失敗しました",
			"project_name": projectName,
			"project_memo": projectMemo,
		})
	}
}

