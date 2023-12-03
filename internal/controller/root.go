package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
)


type RootController struct {}


func NewRootController() *RootController {
	return &RootController{}
}


//GET /
func (ctr *RootController) IndexPage(c *gin.Context) {
	c.Redirect(303, fmt.Sprintf("/%s", jwt.GetUsername(c)))
}