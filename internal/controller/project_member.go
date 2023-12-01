package controller
/*
import (
	"strconv"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type ProjectMemberController struct {
	projectMemberService  service.ProjectMemberService
	projectService  service.ProjectService
}


func NewProjectMemberController() *ProjectMemberController {
	projectMemberService  := service.NewProjectMemberService()
	projectService  := service.NewProjectService()
	return &ProjectMemberController{projectMemberService , projectService}
}


//POST /projects/requests/join
func (ctr *ProjectMemberController) JoinRequest(c *gin.Context) {
	userId := jwt.GetUserId(c)
	p := ctr.projectService .GetProjectByCd(c.PostForm("project_cd"))

	ctr.projectMemberService .JoinRequest(userId, p.ProjectId)
	
	c.Redirect(303, "/projects")
}


//POST /projects/requests/cancel
func (ctr *ProjectMemberController) CancelJoinRequest(c *gin.Context) {
	userId := jwt.GetUserId(c)
	p := ctr.projectService .GetProjectByCd(c.PostForm("project_cd"))

	ctr.projectMemberService .CancelJoinRequest(userId, p.ProjectId)
	
	c.Redirect(303, "/projects")
}


//POST /projects/requests/permit
func (ctr *ProjectMemberController) PermitJoinRequest(c *gin.Context) {
	userId := jwt.GetUserId(c)
	targetUserId, err := strconv.Atoi(c.PostForm("user_id"))
	projectId := ctr.projectService .GetProjectId(userId, c.PostForm("project_cd"))

	if err == nil && projectId != service.GET_PROJECT_ID_NOT_FOUND_INT {
		ctr.projectMemberService .PermitJoinRequest(targetUserId, projectId)
	}

	c.Redirect(303, "/projects/requests")
}


//GET /projects/requests
func (ctr *ProjectMemberController) RequestsPage(c *gin.Context) {
	userId := jwt.GetUserId(c)

	joinrequests, _ := ctr.projectMemberService .GetJoinRequests(userId)
	
	c.HTML(200, "requests.html", gin.H{
		"commons": constant.Commons,
		"joinrequests": joinrequests,
	})
}
*/