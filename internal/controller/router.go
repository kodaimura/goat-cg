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
        rc := newRootController()
        a.GET("/", rc.indexPage)

        pc := newProjectController()
        a.GET("/projects", pc.projectsPage)
        a.GET("/projects/new", pc.createProjectPage)
        a.POST("/projects", pc.createProject)
        a.GET("/projects/:project_id/join", pc.joinRequest)
        a.GET("/projects/:project_id/cancel", pc.cancelJoinRequest)

        ap := a.Group("/:project_cd")
        {
            tc := newTableController()
            ap.GET("/tables", tc.tablesPage)
            ap.GET("/tables/new", tc.createTablePage)
            ap.POST("/tables/new", tc.createTable)
            ap.GET("/tables/:table_id", tc.updateTablePage)
            ap.POST("/tables/:table_id", tc.updateTable)
            ap.DELETE("/tables/:table_id", tc.deleteTable)

            aptt := ap.Group("/tables/:table_id")
            {
                cc := newColumnController()
                aptt.GET("/columns", cc.columnsPage)
                aptt.GET("/columns/new", cc.createColumnPage)
                aptt.POST("/columns/new", cc.createColumn)
                aptt.GET("/columns/:column_id", cc.updateColumnPage)
                aptt.POST("/columns/:column_id", cc.updateColumn)
                aptt.DELETE("/columns/:column_id", cc.deleteColumn)
            }
        }
    }

    api := r.Group("/api", jwt.JwtAuthMiddleware())
    {
        api.GET("/profile", uc.getProfile)
    }
}