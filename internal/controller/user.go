package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goat-cg/config"
	"goat-cg/internal/core/jwt"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/service"
)


type UserController struct {
	userService service.UserService
}


func NewUserController() *UserController {
	userService := service.NewUserService()
	return &UserController{userService}
}


//GET /signup
func (ctr *UserController) SignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{
		"commons": constant.Commons,
	})
}

//GET /login
func (ctr *UserController) LoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"commons": constant.Commons,
	})
}


//POST /signup
func (ctr *UserController) Signup(c *gin.Context) {
	name := c.PostForm("user_name")
	pass := c.PostForm("password")

	result := ctr.userService.Signup(name, pass)

	if result == service.SIGNUP_SUCCESS_INT {
		c.Redirect(303, "/login")

	} else if result == service.SIGNUP_CONFLICT_INT {
		c.HTML(409, "signup.html", gin.H{
			"commons": constant.Commons,
			"error": "ユーザ名が既に使われています。",
		})

	} else {
		c.HTML(500, "signup.html", gin.H{
			"commons": constant.Commons,
			"error": "登録に失敗しました。",
		})
	}
}


//POST /login
func (ctr *UserController) Login(c *gin.Context) {
	name := c.PostForm("user_name")
	pass := c.PostForm("password")

	userId := ctr.userService.Login(name, pass)

	if userId == service.LOGIN_FAILURE_INT {
		c.HTML(401, "login.html", gin.H{
			"commons": constant.Commons,
			"error": "ユーザ名またはパスワードが異なります。",
		})
		c.Abort()
		return
	}

	jwtStr := ctr.userService.GenerateJWT(userId)

	if jwtStr == service.GENERATE_JWT_FAILURE_STR {
		c.HTML(500, "login.html", gin.H{
			"commons": constant.Commons,
			"error": "ログインに失敗しました。",
		})
		c.Abort()
		return
	}

	cf := config.GetConfig()
	c.SetCookie(jwt.COOKIE_KEY_JWT, jwtStr, int(jwt.JWT_EXPIRES), "/", cf.AppHost, false, true)
	c.Redirect(303, "/")
}


//GET /logout
func (ctr *UserController) Logout(c *gin.Context) {
	cf := config.GetConfig()
	c.SetCookie(jwt.COOKIE_KEY_JWT, "", 0, "/", cf.AppHost, false, true)
	c.Redirect(303, "/login")
}


//GET /api/profile
func (ctr *UserController) GetProfile(c *gin.Context) {
	user, err := ctr.userService.GetProfile(jwt.GetUserId(c))

	if err != nil {
		c.JSON(500, gin.H{"error": http.StatusText(500)})
		c.Abort()
		return
	}

	c.JSON(200, user)
}