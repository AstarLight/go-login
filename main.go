package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)


func main() {

	r := gin.Default()

	//no login
	r.Use(CommonRateLimit())
	r.Use(CommonBlacklist())

	r.GET("/home_page", HomePage)
	r.GET("/login_page", LoginPage)

	r.POST("/sign_in", SignIn)
	r.POST("/sign_up", SignUp)
	r.POST("/forget_passwd", ForgetPasswd)


	// needlogin
	needlogin := r.Group("/user")
	needlogin.Use(NeedLogin())
	needlogin.Use(UIDRateLimit())
	needlogin.Use(UIDBlacklist())
	{
		// 管理员接口
		needlogin.POST("/update_passwd", UpdatePasswd)
		needlogin.POST("/sign_out", SignOut)

	}

	fmt.Printf("server run, listen port %s", 9999)
	r.Run(":9999")

}
