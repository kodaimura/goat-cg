package controller

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type projectUserController struct {
	projectUserService  service.ProjectUserService
	projectService  service.ProjectService
}


func newProjectUserController() *projectUserController {
	projectUserService  := service.NewProjectUserService()
	projectService  := service.NewProjectService()
	return &projectUserController{projectUserService , projectService}
}


//POST /projects/requests/join
func (ctr *projectUserController) joinRequest(c *gin.Context) {
	userId := jwt.GetUserId(c)
	p := ctr.projectService .GetProjectByCd(c.PostForm("project_cd"))

	ctr.projectUserService .JoinRequest(userId, p.ProjectId)
	
	c.Redirect(303, "/projects")
}


//POST /projects/requests/cancel
func (ctr *projectUserController) cancelJoinRequest(c *gin.Context) {
	userId := jwt.GetUserId(c)
	p := ctr.projectService .GetProjectByCd(c.PostForm("project_cd"))

	ctr.projectUserService .CancelJoinRequest(userId, p.ProjectId)
	
	c.Redirect(303, "/projects")
}


//POST /projects/requests/permit
func (ctr *projectUserController) permitJoinRequest(c *gin.Context) {
	userId := jwt.GetUserId(c)
	targetUserId, err := strconv.Atoi(c.PostForm("user_id"))
	projectId := ctr.projectService .GetProjectId(userId, c.PostForm("project_cd"))

	if err == nil && projectId != service.GET_PROJECT_ID_NOT_FOUND_INT {
		ctr.projectUserService .PermitJoinRequest(targetUserId, projectId)
	}

	c.Redirect(303, "/projects/requests")
}


//GET /projects/requests
func (ctr *projectUserController) requestsPage(c *gin.Context) {
	userId := jwt.GetUserId(c)

	joinrequests, _ := ctr.projectUserService .GetJoinRequests(userId)
	
	c.HTML(200, "requests.html", gin.H{
		"commons": constant.Commons,
		"joinrequests": joinrequests,
	})
}
