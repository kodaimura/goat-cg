package controller

import (
    "github.com/gin-gonic/gin"
    
    "goat-cg/internal/core/jwt"
    "goat-cg/internal/constant"
)


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