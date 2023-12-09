package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/errs"
	"goat-cg/internal/service"
	"goat-cg/internal/model"
)


type MemberController struct {
	memberService service.MemberService
	projectService service.ProjectService
}


func NewMemberController() *MemberController {
	memberService := service.NewMemberService()
	projectService := service.NewProjectService()
	return &MemberController{memberService , projectService}
}


//GET /:username/:project_name/members
func (mc *MemberController) MemberPage (c *gin.Context) {
	project := c.Keys["project"].(model.Project)

	members, _ := mc.memberService.GetMembers(project.ProjectId)

	c.HTML(200, "members.html", gin.H{
		"project": project, 
		"members": members,
	})
}


//POST /:username/:project_name/members/invite
func (mc *MemberController) Invite (c *gin.Context) {
	project := c.Keys["project"].(model.Project)
	email := c.PostForm("email")

	if email == jwt.GetEmail(c) {
		members, _ := mc.memberService.GetMembers(project.ProjectId)
		c.HTML(400, "members.html", gin.H{
			"project": project, 
			"members": members,
			"email": email,
			"error": "Cannot invite yourself.",
		})
	}

	err := mc.memberService.Invite(project.ProjectId, email)

	if err == nil {
		c.Redirect(303, fmt.Sprintf("/%s/%s/members", c.Param("username"), c.Param("project_name")))
		return
	} 
	
	members, _ := mc.memberService.GetMembers(project.ProjectId)

	if _, ok := err.(errs.NotFoundError); ok {
		c.HTML(400, "members.html", gin.H{
			"project": project, 
			"members": members,
			"email": email,
			"error": "User not found.",
		})

	} else if _, ok := err.(errs.AlreadyRegisteredError); ok{
		c.HTML(409, "members.html", gin.H{
			"project": project, 
			"members": members,
			"email": email,
			"error": "This user has already been invited.",
		})
	} else {
		c.HTML(500, "members.html", gin.H{
			"project": project, 
			"members": members,
			"email": email,
			"error": "invitation failed.",
		})
	}
}