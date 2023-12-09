package jwt

import (
	"log"
	"github.com/gin-gonic/gin"
)

/*
JwtPayload拡張
*/

type CustomClaims struct {
	UserId int
	Username string
	Email string
}


func GetUserId (c *gin.Context) int {
	pl := c.Keys[CONTEXT_KEY_PAYLOAD]
	if pl == nil {
		log.Panic("Error: GetUserId")
		return -1
	} else {
		return pl.(JwtPayload).UserId
	}	
}

func GetUsername (c *gin.Context) string {
	pl := c.Keys[CONTEXT_KEY_PAYLOAD]
	if pl == nil {
		log.Panic("Error: GetUsername")
		return ""
	} else {
		return pl.(JwtPayload).Username
	}
}

func GetEmail (c *gin.Context) string {
	pl := c.Keys[CONTEXT_KEY_PAYLOAD]
	if pl == nil {
		log.Panic("Error: GetEmail")
		return ""
	} else {
		return pl.(JwtPayload).Email
	}
}