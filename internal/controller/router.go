package controller

import (
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
)


func SetRouter(r *gin.Engine) {
    uc := newUserController()
    r.GET("/signup", uc.signupPage)
    r.GET("/login", uc.loginPage)
    r.GET("/logout", uc.logout)

    r.POST("/signup", uc.signup)
    r.POST("/login", uc.login)

    a := r.Group("/", jwt.JwtAuthMiddleware())
    {
        pc := newProjectController()
        a.GET("/projects", pc.projectsPage)
        a.GET("/projects/new", pc.projectPage)
        a.POST("/projects", pc.postProjects)

        ap := a.Group("/:project_cd")
        {
            tc := newTableController()
            ap.GET("/tables", tc.tablesPage)
            ap.POST("/tables", tc.postTable)

            aptt := ap.Group("/tables/:table_id")
            {
                cc := newColumnController()
                aptt.GET("/columns", cc.columnsPage)
                aptt.POST("/columns", cc.postColumns)
            }
        }
    }

    api := r.Group("/api", jwt.JwtAuthMiddleware())
    {
        api.GET("/profile", uc.getProfile)
    }
}