package controller

import (
    "github.com/gin-gonic/gin"
    
    "goat-cg/internal/core/jwt"
    "goat-cg/internal/constant"
)


func setRootRoute(r *gin.Engine) {
    auth := r.Group("/", jwt.JwtAuthMiddleware())

    rc := newRootController()
    auth.GET("/", rc.indexPage)
}


type rootController struct {}


func newRootController() *rootController {
    return &rootController{}
}


//GET /
func (ctr *rootController) indexPage(c *gin.Context) {
    username := jwt.GetUserName(c)
    
    c.HTML(200, "index.html", gin.H{
        "commons": constant.Commons,
        "username": username,
    })
}