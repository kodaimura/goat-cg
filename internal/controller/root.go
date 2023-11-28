package controller

import (
	"github.com/gin-gonic/gin"
)


type RootController struct {}


func NewRootController() *RootController {
	return &RootController{}
}


//GET /
func (ctr *RootController) indexPage(c *gin.Context) {
	c.Redirect(303, "/projects")
}