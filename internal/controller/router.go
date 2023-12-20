package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/internal/core/jwt"
	"goat-cg/internal/middleware"
)


func SetRouter(r *gin.Engine) {
	uc := NewUserController()

	//render HTML or redirect
	r.GET("/signup", uc.SignupPage)
	r.POST("/signup", uc.Signup)
	r.GET("/login", uc.LoginPage)
	r.POST("/login", uc.Login)
	r.GET("/logout", uc.Logout)

	//render HTML or redirect (Authorized request)
	a := r.Group("/", jwt.JwtAuthMiddleware())
	{
		rc := NewRootController()
		
		a.GET("/", rc.IndexPage)

		pc := NewProjectController()

		au := a.Group("/:username")
		{
			au.GET("", pc.ProjectsPage)
			au.GET("/projects", pc.ProjectsPage)
			au.GET("/projects/new", pc.CreateProjectPage)
			au.POST("/projects/new", pc.CreateProject)
			au.GET("/projects/:project_id", pc.UpdateProjectPage)
			au.POST("/projects/:project_id", pc.UpdateProject)
			au.DELETE("/projects/:project_id", pc.DeleteProject)
			au.GET("/account", uc.AccountPage)
			au.POST("/account/password", uc.UpdatePassword)
			au.POST("/account/email", uc.UpdateEmail)

			aup := au.Group("/:project_name", middleware.PathParameterValidationMiddleware())
			{
				tc := NewTableController()
				mc := NewMemberController()

				aup.GET("", tc.TablesPage)
				aup.GET("/tables", tc.TablesPage)
				aup.GET("/tables/new", tc.CreateTablePage)
				aup.POST("/tables/new", tc.CreateTable)
				aup.GET("/tables/:table_id", tc.UpdateTablePage)
				aup.POST("/tables/:table_id", tc.UpdateTable)
				aup.DELETE("/tables/:table_id", tc.DeleteTable)
				aup.GET("/tables/:table_id/log", tc.TableLogPage)
				aup.GET("/members", mc.MembersPage)
				aup.GET("/members/:user_id", mc.MemberPage)
				aup.DELETE("/members/:user_id", mc.DeleteMember)
				aup.POST("/members/invite", mc.Invite)
	
	
				cgc := NewCodegenController()
	
				aup.GET("/codegen", cgc.CodegenPage)
				aup.POST("/codegen/goat", cgc.CodegenGOAT)
	
				aupt := aup.Group("/tables/:table_id")
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

	//response JSON (Authorized request)
	api := r.Group("/api", jwt.JwtAuthApiMiddleware())
	{
		api.GET("/account/profile", uc.GetProfile)
		//api.PUT("/account/username", uc.UpdateUsername)
		//api.PUT("/account/password", uc.UpdatePassword)
		//api.PUT("/account/email", uc.UpdateEmail)
		api.DELETE("/account", uc.DeleteAccount)
	}
}