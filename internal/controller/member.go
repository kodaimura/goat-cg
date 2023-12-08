package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/service"
	"goat-cg/internal/model"
)


type MemberController struct {
	memberService  service.MemberService
	projectService  service.ProjectService
}


func NewMemberController() *MemberController {
	memberService  := service.NewMemberService()
	projectService  := service.NewProjectService()
	return &MemberController{memberService , projectService}
}


//GET /:username/:project_name/members
func (mc *MemberController) MemberPage (c *gin.Context) {
	project := c.Keys["project"].(model.Project)

	c.HTML(200, "members.html", gin.H{
		"project": project, 
	})
}