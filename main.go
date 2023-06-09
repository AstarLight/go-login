package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	// 配置初始化
	ConfInit()
	// 数据库初始化
	DbInit()
	//seesion初始化
	SessionInit()
	// 分布式限流器
	RatelimitInit()

	r := gin.Default()
	r.LoadHTMLGlob("./assets/template/*.html")
	r.Static("/assets/bootstrap", "./assets/bootstrap")

	//no login
	r.GET("/admin_login.html", GetTemplate)
	r.GET("/admin_regist.html", GetTemplate)
	r.GET("/login_regist.html", GetTemplate)

	needLimit:= r.Group("/sign")
	needLimit.Use(CommonBlacklist()) // 黑名单
	needLimit.Use(CommonRateLimit()) // 频率控制
	{
		needLimit.POST("/sign_in", SignIn) // 登录
		needLimit.POST("/sign_up", SignUp) // 注册
	}

	// needlogin 以下接口需要登录态才可访问
	needlogin := r.Group("/user")
	needlogin.Use(NeedLogin())
	needlogin.Use(UserBlacklist())
	needlogin.Use(UserRateLimit())
	{
		needlogin.GET("/update_password_page.html", GetTemplate) // 更新密码页
		needlogin.GET("/home.html", GetTemplate)                 // 用户首页
		needlogin.POST("/update_passwd", UpdatePasswd)           // 更新密码
		needlogin.POST("/sign_out", SignOut)                     // 登出

	}

	sPort := fmt.Sprintf("%d", Conf.Common.ListenPort)
	fmt.Printf("server run, listen port %s", sPort)
	r.Run(":" + sPort)

}
