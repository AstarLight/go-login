package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	//"net/http"
)

//必须登录的请求，从session读user写入context
func NeedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, sess := GetUserFromSession(c)
		if err == nil && sess != nil {
			fmt.Println("NeedLogin in login, ", sess)
			c.Set("USER", sess)
			sess.Store() // 已登录，每次请求都会续期
			c.Next()
			return

		} else {
			// 未登录
			WriteResponseWithCode(c, "未登录", nil, 429)
			//c.Redirect(http.StatusFound, Conf.Common.EnterPage)
			c.Abort()
			return
		}
	}
}

//限流，粒度分为IP,MAC-ADDR
func CommonRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		return
	}
}

//黑名单,限流，粒度分为IP,MAC-ADDR
func CommonBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		return
	}
}

//限流，粒度分为UID
func UIDRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		return
	}
}

//黑名单,限流，粒度分为UID
func UIDBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		return
	}
}
