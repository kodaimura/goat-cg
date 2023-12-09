package controller

import (
	"github.com/gin-gonic/gin"

	"goat-cg/config"
	"goat-cg/internal/core/jwt"
	"goat-cg/internal/core/errs"
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
func (uc *UserController) SignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{})
}

//GET /login
func (uc *UserController) LoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

//GET /account
func (uc *UserController) AccountPage(c *gin.Context) {
	user, _ := uc.userService.GetProfile(jwt.GetUserId(c))
	c.HTML(200, "account.html", gin.H{
		"username": user.Username,
		"email": user.Email,
	})
}

//POST /signup
func (uc *UserController) Signup(c *gin.Context) {
	name := c.PostForm("username")
	pass := c.PostForm("password")
	email := c.PostForm("email")

	err := uc.userService.Signup(name, pass, email)

	if err != nil {
		if _, ok := err.(errs.UniqueConstraintError); ok {
			if err.(errs.UniqueConstraintError).Column == "username" {
				c.HTML(409, "signup.html", gin.H{
					"username": name,
					"password": pass,
					"email": email,
					"error": "This Username is already taken.",
				})
			} else {
				c.HTML(409, "signup.html", gin.H{
					"username": name,
					"password": pass,
					"email": email,
					"error": "This Email is already taken.",
				})
			}
			
		} else {
			c.HTML(500, "signup.html", gin.H{
				"username": name,
				"password": pass,
				"email": email,
				"error": "error occurred.",
			})
		}
		c.Abort()
		return
	}

	c.Redirect(303, "/login")
}


//POST /login
func (uc *UserController) Login(c *gin.Context) {
	name := c.PostForm("username")
	pass := c.PostForm("password")

	user, err := uc.userService.Login(name, pass)

	if err != nil {
		c.HTML(401, "login.html", gin.H{
			"username": name,
			"password": pass,
			"error": "Incorrect Username or Password.",
		})
		c.Abort()
		return
	}

	jwtStr, err := uc.userService.GenerateJWT(user.UserId)

	if err != nil {
		c.HTML(500, "login.html", gin.H{
			"username": name,
			"password": pass,
			"error": "error occurred.",
		})
		c.Abort()
		return
	}

	cf := config.GetConfig()
	c.SetCookie(jwt.COOKIE_KEY_JWT, jwtStr, int(jwt.JWT_EXPIRES), "/", cf.AppHost, false, true)
	c.Redirect(303, "/")
}


//GET /logout
func (uc *UserController) Logout(c *gin.Context) {
	cf := config.GetConfig()
	c.SetCookie(jwt.COOKIE_KEY_JWT, "", 0, "/", cf.AppHost, false, true)
	c.Redirect(303, "/login")
}


//POST /account/password
func (uc *UserController) UpdatePassword(c *gin.Context) {
	id := jwt.GetUserId(c)
	pass := c.PostForm("password")
	newPass := c.PostForm("new_password")

	user, err := uc.userService.Login(jwt.GetUsername(c), pass)

	if err != nil {
		user, _ = uc.userService.GetProfile(jwt.GetUserId(c))
		c.HTML(400, "account.html", gin.H{
			"password_error": "Incorrect Current Password.",
			"username": user.Username,
			"email": user.Email,
		})
		c.Abort()
		return
	}

	if uc.userService.UpdatePassword(id, newPass) != nil {
		c.HTML(500, "account.html", gin.H{
			"password_error": "error occurred.",
			"username": user.Username,
			"email": user.Email,
		})
		c.Abort()
		return
	}

	c.Redirect(303, "/logout")
}


//POST /account/email
func (uc *UserController) UpdateEmail(c *gin.Context) {
	id := jwt.GetUserId(c)
	email := c.PostForm("email")
	
	err := uc.userService.UpdateEmail(id, email)
	if err != nil {
		if _, ok := err.(errs.UniqueConstraintError); ok {
			c.HTML(409, "account.html", gin.H{
				"email_error": "This Email is already taken.",
				"username": jwt.GetUsername(c),
				"email": email,
			})
		} else {
			c.HTML(500, "account.html", gin.H{
				"email_error": "error occurred.",
				"username": jwt.GetUsername(c),
				"email": email,
			})
		}
		c.Abort()
		return
	}

	c.Redirect(303, "/logout")
}


//GET /api/account/profile
func (uc *UserController) GetProfile(c *gin.Context) {
	user, err := uc.userService.GetProfile(jwt.GetUserId(c))

	if err != nil {
		c.JSON(500, gin.H{})
		c.Abort()
		return
	}

	c.JSON(200, user)
}


//DELETE /api/account
func (uc *UserController) DeleteAccount(c *gin.Context) {
	id := jwt.GetUserId(c)

	if uc.userService.DeleteUser(id) != nil {
		c.JSON(500, gin.H{"error": "error occurred."})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{})
}