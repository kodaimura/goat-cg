package controller

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/service"
	"goat-cg/internal/model"
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
		return
	}

	var project model.Project
	project.ProjectName = projectName
	project.ProjectMemo= projectMemo
	
	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "project.html", gin.H{
			"username": username,
			"error": "ProjectName must be unique.",
			"project": project,
		})

	} else {
		c.HTML(500, "project.html", gin.H{
			"username": username,
			"error": "error occurred.",
			"project": project,
		})
	}
}


//GET /:username/projects/:project_id
func (ctr *ProjectController) UpdateProjectPage(c *gin.Context) {
	userId := jwt.GetUserId(c)
	username := jwt.GetUsername(c)
	projectId, err := strconv.Atoi(c.Param("project_id"))

	if err != nil || c.Param("username") != username {
		c.HTML(404, "404error.html", gin.H{})
		c.Abort()
		return
	}

	project, err := ctr.projectService.GetProject(projectId)

	if err != nil || project.UserId != userId {
		c.HTML(404, "404error.html", gin.H{})
		c.Abort()
		return
	}

	c.HTML(200, "project.html", gin.H{
		"username": username,
		"project": project, 
	})
}


//POST /:username/projects/:project_id
func (ctr *ProjectController) UpdateProject(c *gin.Context) {
	username := jwt.GetUsername(c)
	projectId, err := strconv.Atoi(c.Param("project_id"))

	if err != nil || c.Param("username") != username {
		c.HTML(404, "404error.html", gin.H{})
		c.Abort()
		return
	}

	projectName := c.PostForm("project_name")
	projectMemo := c.PostForm("project_memo")

	err = ctr.projectService.UpdateProject(username, projectId, projectName, projectMemo)
	if err == nil {
		c.Redirect(303, fmt.Sprintf("/%s", username))
		return
	} 
	
	var project model.Project
	project.ProjectName = projectName
	project.ProjectMemo= projectMemo
	project.Username= username

	if _, ok := err.(errs.UniqueConstraintError); ok {
		c.HTML(409, "project.html", gin.H{
			"username": username,
			"error": "ProjectName must be unique.",
			"project": project,
		})

	} else {
		c.HTML(500, "project.html", gin.H{
			"username": username,
			"error": "error occurred.",
			"project": project,
		})
	}
}


//DELETE /:username/projects/:project_id
func (ctr *ProjectController) DeleteProject(c *gin.Context) {
	username := jwt.GetUsername(c)
	projectId, err := strconv.Atoi(c.Param("project_id"))

	if err != nil || c.Param("username") != username {
		c.HTML(404, "404error.html", gin.H{})
		c.Abort()
		return
	}
	ctr.projectService.DeleteProject(projectId)

	c.Redirect(303, fmt.Sprintf("/%s/%s/tables", c.Param("username"), c.Param("project_name")))

}
