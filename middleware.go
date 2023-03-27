package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

func IsBlackIP(IP string) bool {
	for _, ip := range Conf.Blacklist.IPList {
		if ip == IP {
			return true
		}
	}

	return false
}

func IsBlackUsername(username string) bool {
	for _, uname := range Conf.Blacklist.UsernameList {
		if username == uname {
			return true
		}
	}

	return false
}


//限流，粒度分为IP
func CommonRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if IsReachLimit(clientIP, "ip") {
			WriteResponseWithCode(c, "访问次数过多", nil, http.StatusTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

//黑名单,限流，粒度分为IP
func CommonBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if IsBlackIP(clientIP) {
			WriteResponseWithCode(c, "禁止访问", nil, 403)
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

//限流，粒度分为UID
func UserRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := GetCtxUser(c)
		if sess != nil {
			username := sess.Username

			if IsReachLimit(username, "user") {
				WriteResponseWithCode(c, "访问次数过多", nil, http.StatusTooManyRequests)
				c.Abort()
				return
			}
		}
		c.Next()
		return
	}
}

//黑名单,限流，粒度分为UID
func UserBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := GetCtxUser(c)
		if sess != nil {
			if IsBlackUsername(sess.Username) {
				WriteResponseWithCode(c, "禁止访问", nil, 403)
				c.Abort()
				return
			}
		}
		c.Next()
		return
	}
}
