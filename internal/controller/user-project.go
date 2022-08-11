package controller

import (
    "strconv"
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
    "goat-cg/internal/shared/constant"
    "goat-cg/internal/service"
)


type userProjectController struct {
    upServ service.UserProjectService
    pServ service.ProjectService
}


func newUserProjectController() *userProjectController {
    upServ := service.NewUserProjectService()
    pServ := service.NewProjectService()
    return &userProjectController{upServ, pServ}
}


//POST /projects/requests/join
func (ctr *userProjectController) joinRequest(c *gin.Context) {
    userId := jwt.GetUserId(c)
    p := ctr.pServ.GetProjectByCd(c.PostForm("project_cd"))

    ctr.upServ.JoinRequest(userId, p.ProjectId)
    
    c.Redirect(303, "/projects")
}


//POST /projects/requests/cancel
func (ctr *userProjectController) cancelJoinRequest(c *gin.Context) {
    userId := jwt.GetUserId(c)
    p := ctr.pServ.GetProjectByCd(c.PostForm("project_cd"))

    ctr.upServ.CancelJoinRequest(userId, p.ProjectId)
    
    c.Redirect(303, "/projects")
}


//POST /projects/requests/permit
func (ctr *userProjectController) permitJoinRequest(c *gin.Context) {
    userId := jwt.GetUserId(c)
    targetUserId, err := strconv.Atoi(c.PostForm("user_id"))
    projectId := ctr.pServ.GetProjectId(userId, c.PostForm("project_cd"))

    if err != nil || projectId == service.GET_PROJECT_ID_NOT_FOUND_INT {
        c.Redirect(303, "/projects/requests")
        return
    }

    ctr.upServ.PermitJoinRequest(targetUserId, projectId)
    
    c.Redirect(303, "/projects/requests")
}


//GET /projects/requests
func (ctr *userProjectController) requestsPage(c *gin.Context) {
    userId := jwt.GetUserId(c)

    joinrequests, _ := ctr.upServ.GetJoinRequests(userId)
    
    c.HTML(200, "requests.html", gin.H{
        "commons": constant.Commons,
        "joinrequests": joinrequests,
    })
}
