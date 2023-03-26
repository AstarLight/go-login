package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// 配置初始化
	ConfInit()
	// 数据库初始化
	DbInit()
	//redis初始化
	RedisInit()

	r := gin.Default()
	r.LoadHTMLGlob("template/*.html")
	r.Static("/static", "./static")

	//no login
	r.Use(CommonRateLimit()) // 频率控制
	r.Use(CommonBlacklist()) // 黑名单

	r.GET("/enter.html", GetTemplate)
	r.GET("/home.html", GetTemplate)
	r.GET("/login.html", GetTemplate)
	r.GET("/regist.html", GetTemplate)

	r.POST("/sign_in", SignIn)             // 登录
	r.POST("/sign_up", SignUp)             // 注册
	r.POST("/forget_passwd", ForgetPasswd) // 忘记密码

	// needlogin 以下接口需要登录态才可访问
	needlogin := r.Group("/user")
	needlogin.Use(NeedLogin())
	needlogin.Use(UIDRateLimit())
	needlogin.Use(UIDBlacklist())
	{
		needlogin.POST("/update_passwd", UpdatePasswd) // 更新密码
		needlogin.POST("/sign_out", SignOut)           // 登出

		needlogin.GET("/update_password_page.html", GetTemplate)

	}

	sPort := fmt.Sprintf("%d", Conf.Common.ListenPort)
	fmt.Printf("server run, listen port %s", sPort)
	r.Run(":" + sPort)

}
