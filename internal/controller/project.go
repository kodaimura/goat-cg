package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
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


//GET /projects
func (ctr *ProjectController) ProjectsPage(c *gin.Context) {
	userId := jwt.GetUserId(c)
	username := jwt.GetUsername(c)

	projects, _ := ctr.projectService.GetProjects(userId)
	member_projects, _ := ctr.projectService.GetMemberProjects(userId)

	c.HTML(200, "projects.html", gin.H{
		"commons": constant.Commons,
		"username": username,
		"projects": projects,
		"member_projects": member_projects,
	})
}


//GET /projects/new
func (ctr *ProjectController) CreateProjectPage(c *gin.Context) {
	
	c.HTML(200, "project.html", gin.H{
		"commons": constant.Commons,
	})
}


//POST /projects
func (ctr *ProjectController) CreateProject(c *gin.Context) {
	projectName := c.PostForm("project_name")
	projectMemo := c.PostForm("project_memo")
	result := ctr.projectService.CreateProject(jwt.GetUserId(c), jwt.GetUsername(c), projectName, projectMemo)
	
	if result == service.CREATE_PROJECT_SUCCESS_INT {
		c.Redirect(303, "/projects")
		return
	}

	if result == service.CREATE_PROJECT_CONFLICT_INT {
		c.HTML(409, "project.html", gin.H{
			"commons": constant.Commons,
			"error": "プロジェクト名が重複して使われています",
			"project_name": projectName,
			"project_memo": projectMemo,
		})
	} else {
		c.HTML(500, "project.html", gin.H{
			"commons": constant.Commons,
			"error": "登録に失敗しました",
			"project_name": projectName,
			"project_memo": projectMemo,
		})

	}
}

