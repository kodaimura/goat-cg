package controller

import (
    "github.com/gin-gonic/gin"

    "goat-cg/internal/core/jwt"
)


func SetRouter(r *gin.Engine) {
    uc := newUserController()

    //render HTML or redirect
    r.GET("/signup", uc.signupPage)
    r.GET("/login", uc.loginPage)
    r.GET("/logout", uc.logout)
    r.POST("/signup", uc.signup)
    r.POST("/login", uc.login)

    //render HTML or redirect (Authorized request)
    a := r.Group("/", jwt.JwtAuthMiddleware())
    {
        rc := newRootController()
        
        a.GET("/", rc.indexPage)

        pc := newProjectController()

        a.GET("/projects", pc.projectsPage)
        a.GET("/projects/new", pc.createProjectPage)
        a.POST("/projects", pc.createProject)

        upc := newProjectUserController()

        a.GET("/projects/requests", upc.requestsPage)
        a.POST("/projects/requests/join", upc.joinRequest)
        a.POST("/projects/requests/cancel", upc.cancelJoinRequest)
        a.POST("/projects/requests/permit", upc.permitJoinRequest)


        ap := a.Group("/:project_cd")
        {
            tc := newTableController()

            ap.GET("/tables", tc.tablesPage)
            ap.GET("/tables/new", tc.createTablePage)
            ap.POST("/tables/new", tc.createTable)
            ap.GET("/tables/:table_id", tc.updateTablePage)
            ap.POST("/tables/:table_id", tc.updateTable)
            ap.DELETE("/tables/:table_id", tc.deleteTable)
            ap.GET("/tables/:table_id/log", tc.tableLogPage)


            cgc := newCodegenController()

            ap.GET("/codegen", cgc.codegenPage)
            ap.POST("/codegen/goat", cgc.codegenGOAT)
            ap.POST("/codegen/ddl", cgc.codegenDDL)

            aptt := ap.Group("/tables/:table_id")
            {
                cc := newColumnController()

                aptt.GET("/columns", cc.columnsPage)
                aptt.GET("/columns/new", cc.createColumnPage)
                aptt.POST("/columns/new", cc.createColumn)
                aptt.GET("/columns/:column_id", cc.updateColumnPage)
                aptt.POST("/columns/:column_id", cc.updateColumn)
                aptt.DELETE("/columns/:column_id", cc.deleteColumn)
                aptt.GET("/columns/:column_id/log", cc.columnLogPage)
            }
        }
    }

    //response JSON
    api := r.Group("/api")
    {
        uac := newUserApiController()

        api.POST("/signup", uac.signup)
        api.POST("/login", uac.login)
        api.GET("/logout", uac.logout)


        //response JSON (Authorized request)
        a := api.Group("/", jwt.JwtAuthApiMiddleware())
        {
            a.GET("/profile", uac.getProfile)
            a.PUT("/username", uac.changeUsername)
            a.POST("/username", uac.changeUsername)
            a.PUT("/password", uac.changePassword)
            a.POST("/password", uac.changePassword)
            a.DELETE("/account", uac.deleteUser)
        }
    }
}