package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
)


type RootController struct {}


func NewRootController() *RootController {
	return &RootController{}
}


//GET /
func (ctr *RootController) IndexPage(c *gin.Context) {
	username := jwt.GetUsername(c)
	c.Redirect(303, fmt.Sprintf("/%s/projects", c.Param(username)))
}