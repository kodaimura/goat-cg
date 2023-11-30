package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/middleware/jwt"
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

		au := a.Group("/:username")
		{
			au.GET("/projects", pc.ProjectsPage)
			au.GET("/projects/new", pc.CreateProjectPage)
			au.POST("/projects", pc.CreateProject)

			aup := a.Group("/:project_name", middleware.PathParameterValidationMiddleware())
			{
				tc := NewTableController()
	
				aup.GET("/tables", tc.TablesPage)
				aup.GET("/tables/new", tc.CreateTablePage)
				aup.POST("/tables/new", tc.CreateTable)
				aup.GET("/tables/:table_id", tc.UpdateTablePage)
				aup.POST("/tables/:table_id", tc.UpdateTable)
				aup.DELETE("/tables/:table_id", tc.DeleteTable)
				aup.GET("/tables/:table_id/log", tc.TableLogPage)
	
	
				cgc := NewCodegenController()
	
				aup.GET("/codegen", cgc.CodegenPage)
				aup.POST("/codegen/goat", cgc.CodegenGOAT)
				aup.POST("/codegen/ddl", cgc.CodegenDDL)
	
				aupt := ap.Group("/tables/:table_id")
				{
					cc := NewColumnController()
	
					aupt.GET("/columns", cc.ColumnsPage)
					aupt.GET("/columns/new", cc.CreateColumnPage)
					aupt.POST("/columns/new", cc.CreateColumn)
					aupt.GET("/columns/:column_id", cc.UpdateColumnPage)
					aupt.POST("/columns/:column_id", cc.UpdateColumn)
					aupt.DELETE("/columns/:column_id", cc.DeleteColumn)
					aupt.GET("/columns/:column_id/log", cc.ColumnLogPage)
				}
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