package controller

import (
    "github.com/gin-gonic/gin"
)


func SetRouter(r *gin.Engine) {
    setUserRoute(r)
    setRootRoute(r)
    setProjectRoute(r)
    setTableRoute(r)
} 