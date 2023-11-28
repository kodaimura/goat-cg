package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
)


func SetRouter(r *gin.Engine) {
	uc := NewUserController()

	//render HTML or redirect
	r.GET("/signup", uc.SignupPage)
	r.GET("/login", uc.LoginPage)
	r.GET("/logout", uc.Logout)
	r.POST("/signup", uc.Signup)
	r.POST("/login", uc.Login)

	//render HTML or redirect (Authorized request)
	a := r.Group("/", jwt.JwtAuthMiddleware())
	{
		rc := NewRootController()
		
		a.GET("/", rc.IndexPage)

		pc := NewProjectController()

		a.GET("/projects", pc.ProjectsPage)
		a.GET("/projects/new", pc.CreateProjectPage)
		a.POST("/projects", pc.CreateProject)

		upc := NewProjectUserController()

		a.GET("/projects/requests", upc.RequestsPage)
		a.POST("/projects/requests/join", upc.JoinRequest)
		a.POST("/projects/requests/cancel", upc.CancelJoinRequest)
		a.POST("/projects/requests/permit", upc.PermitJoinRequest)


		ap := a.Group("/:project_cd")
		{
			tc := NewTableController()

			ap.GET("/tables", tc.TablesPage)
			ap.GET("/tables/new", tc.CreateTablePage)
			ap.POST("/tables/new", tc.CreateTable)
			ap.GET("/tables/:table_id", tc.UpdateTablePage)
			ap.POST("/tables/:table_id", tc.UpdateTable)
			ap.DELETE("/tables/:table_id", tc.DeleteTable)
			ap.GET("/tables/:table_id/log", tc.TableLogPage)


			cgc := NewCodegenController()

			ap.GET("/codegen", cgc.CodegenPage)
			ap.POST("/codegen/goat", cgc.CodegenGOAT)
			ap.POST("/codegen/ddl", cgc.CodegenDDL)

			aptt := ap.Group("/tables/:table_id")
			{
				cc := NewColumnController()

				aptt.GET("/columns", cc.ColumnsPage)
				aptt.GET("/columns/new", cc.CreateColumnPage)
				aptt.POST("/columns/new", cc.CreateColumn)
				aptt.GET("/columns/:column_id", cc.UpdateColumnPage)
				aptt.POST("/columns/:column_id", cc.UpdateColumn)
				aptt.DELETE("/columns/:column_id", cc.DeleteColumn)
				aptt.GET("/columns/:column_id/log", cc.ColumnLogPage)
			}
		}
	}

	//response JSON
	api := r.Group("/api")
	{
		uac := NewUserApiController()

		api.POST("/signup", uac.Signup)
		api.POST("/login", uac.Login)
		api.GET("/logout", uac.Logout)


		//response JSON (Authorized request)
		a := api.Group("/", jwt.JwtAuthApiMiddleware())
		{
			a.GET("/profile", uac.GetProfile)
			a.PUT("/username", uac.ChangeUsername)
			a.POST("/username", uac.ChangeUsername)
			a.PUT("/password", uac.ChangePassword)
			a.POST("/password", uac.ChangePassword)
			a.DELETE("/account", uac.DeleteUser)
		}
	}
}