package controller

import (
	"github.com/gin-gonic/gin"
)


type rootController struct {}


func newRootController() *rootController {
	return &rootController{}
}


//GET /
func (ctr *rootController) indexPage(c *gin.Context) {
	c.Redirect(303, "/projects")
}