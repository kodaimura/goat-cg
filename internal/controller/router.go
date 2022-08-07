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
        a.GET("/projects/new", pc.createProjectPage)
        a.POST("/projects", pc.createProject)

        ap := a.Group("/:project_cd")
        {
            tc := newTableController()
            ap.GET("/tables", tc.tablesPage)
            ap.GET("/tables/new", tc.createTablePage)
            ap.POST("/tables/new", tc.createTable)
            ap.GET("/tables/:table_id", tc.updateTablePage)
            ap.POST("/tables/:table_id", tc.updateTable)

            aptt := ap.Group("/tables/:table_id")
            {
                cc := newColumnController()
                aptt.GET("/columns", cc.columnsPage)
                aptt.GET("/columns/new", cc.createColumnPage)
                aptt.POST("/columns/new", cc.createColumn)
                aptt.GET("/columns/:column_id", cc.updateColumnPage)
                aptt.POST("/columns/:column_id", cc.updateColumn)
            }
        }
    }

    api := r.Group("/api", jwt.JwtAuthMiddleware())
    {
        api.GET("/profile", uc.getProfile)
    }
}