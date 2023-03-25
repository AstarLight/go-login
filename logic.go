package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"strings"
)

// 登录
func SignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	err := IsValidName(username)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	err = IsValidPasswd(password)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	var user User
	user.Name = username
	has, err := GetUserFromDbByName(&user)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	if !has {
		WriteResponseWithCode(c, "账号不存在", nil, 0)
		return
	}

	if user.Passwd != GenMD5WithSalt(password, user.Salt) {
		WriteResponseWithCode(c, "密码不正确", nil, 0)
		return
	}

	//账密验证通过，生成session
	session := NewSession(&user)
	err = session.Store() // 存储session到redis
	if err != nil {
		WriteResponseWithCode(c, "登录失败，请重试", nil, 0)
		return
	}

	// 更新登录IP和登录时间
	var updates = map[string]interface{}{
		"LastLoginUnix": time.Now().Unix(),
		"LastLoginIp":   c.ClientIP(),
	}
	err = DBUpdateUser(updates)
	if err != nil {
		fmt.Println("DBUpdateUser err: ", err)
	}

	//登录成功，重定向到首页
	c.SetCookie("SESSION", session.ID, 0, "", "", false, true)
	c.SetCookie("USERNAME", session.Username, 0, "", "", false, true)
	c.SetCookie("UID", string(session.UID), 0, "", "", false, true)
	c.Redirect(http.StatusFound, Conf.Common.HomePage) // 重定向跳回首页

}

func GetCtxUser(c *gin.Context) *Session {
	return c.Value("user").(*Session)
}

// 登出
func SignOut(c *gin.Context) {
	sess := GetCtxUser(c)
	if sess == nil {
		WriteResponseWithCode(c, "尚未登录", nil, 0)
		return
	}

	err := sess.Del()
	if err != nil {
		WriteResponseWithCode(c, "注销失败，请重试", nil, 0)
		return
	}

	WriteResponseWithCode(c, "注销成功", nil, 0)
	c.SetCookie("SESSION", "", 0, "", "", false, true)
	c.SetCookie("USERNAME", "", 0, "", "", false, true)
	c.SetCookie("UID", "", 0, "", "", false, true)

	c.Redirect(http.StatusFound, Conf.Common.HomePage)
}

// 注册
func SignUp(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	retryPassword := c.PostForm("retry_password")
	email := c.PostForm("email")

	if retryPassword != password {
		WriteResponseWithCode(c, "两次输入的密码不一致", nil, 0)
		return
	}

	err := IsValidName(username)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	err = IsValidEmail(email)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	err = IsValidPasswd(password)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	if UserExists("username", username) {
		WriteResponseWithCode(c, "用户名已存在", nil, 0)
		return
	}
	if UserExists("email", email) {
		WriteResponseWithCode(c, "邮箱已被注册", nil, 0)
		return
	}

	user := CreateUser(username, password, email)

	// 插入数据库
	err = DBInsertNewUser(user)
	if err != nil {
		WriteResponseWithCode(c, "注册失败，请稍后重新注册", nil, 0)
		return
	}
	WriteResponseWithCode(c, "注册成功，请登录", nil, 0)
	c.Redirect(http.StatusFound, Conf.Common.LoginPage)
}

// 忘记密码
func ForgetPasswd(c *gin.Context) {
	// 给用户邮箱发一封重置密码的邮件，里面有code。用户填入CODE+new passwd

}

// 更新密码
func UpdatePasswd(c *gin.Context) {
	user := GetCtxUser(c)
	if user == nil {
		WriteResponseWithCode(c, "尚未登录", nil, 0)
		return
	}

	password := c.PostForm("password")

	err := IsValidPasswd(password)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	var updates = map[string]interface{}{
		"passwd": password,
	}

	err = DBUpdateUser(updates)
	if err != nil {
		WriteResponseWithCode(c, "修改密码失败，请重试", nil, 0)
		return
	}

	WriteResponseWithCode(c, "修改密码成功", nil, 0)
	c.Redirect(http.StatusFound, Conf.Common.HomePage)
}


func GetTemplate(c *gin.Context) {
	path := c.Request.URL.Path
	arr := strings.Split(path, "/")
	html := arr[len(arr)-1]
	c.HTML(http.StatusOK, html, nil)
}



