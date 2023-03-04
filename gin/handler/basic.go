package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qinchy/hellogo/gin/globalvar"
	"github.com/qinchy/hellogo/gin/types"
	"log"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// StartPage handler
func StartPage(c *gin.Context) {
	var person types.Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(200, "Success")
}

// Getting handler
func Getting(c *gin.Context) {
	// 获取用户，它是由 BasicAuth 中间件设置的
	user := c.MustGet(gin.AuthUserKey).(string)
	if secret, ok := globalvar.Secrets[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}
