package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
	"goat-cg/internal/model/entity"
)


type projectController struct {
	pServ service.ProjectService
}


func newProjectController() *projectController {
	pServ := service.NewProjectService()
	return &projectController{pServ}
}


//GET /projects
func (ctr *projectController) projectsPage(c *gin.Context) {
	userId := jwt.GetUserId(c)
	projectCd := c.Query("project_cd")
	var project entity.Project

	projects, _ := ctr.pServ.GetProjects(userId)
	projects2, _ := ctr.pServ.GetProjectsPendingApproval(userId)
	
	if projectCd != "" {
		project = ctr.pServ.GetProjectByCd(projectCd)
	}

	c.HTML(200, "projects.html", gin.H{
		"commons": constant.Commons,
		"projects": projects,
		"projects2": projects2,
		"project":project,
	})
}


//GET /projects/new
func (ctr *projectController) createProjectPage(c *gin.Context) {
	
	c.HTML(200, "project.html", gin.H{
		"commons": constant.Commons,
	})
}


//POST /projects
func (ctr *projectController) createProject(c *gin.Context) {
	projectCd := c.PostForm("project_cd")
	projectName := c.PostForm("project_name")
	result := ctr.pServ.CreateProject(jwt.GetUserId(c), projectCd, projectName)
	
	if result == service.CREATE_PROJECT_SUCCESS_INT {
		c.Redirect(303, "/projects")
		return
	}

	if result == service.CREATE_PROJECT_CONFLICT_INT {
		c.HTML(409, "project.html", gin.H{
			"commons": constant.Commons,
			"error": "ProjectCd が既に使われています",
			"project_cd": projectCd,
			"project_name": projectName,
		})
	} else {
		c.HTML(500, "project.html", gin.H{
			"commons": constant.Commons,
			"error": "登録に失敗しました",
			"project_cd": projectCd,
			"project_name": projectName,
		})

	}
}

