package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)


func SignIn(c *gin.Context) error {
	username := c.PostForm("username")
	password := c.PostForm("password")

	err := IsValidName(username)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, RESPONSE_OK)
		return
	}

	err = IsValidPasswd(password)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, RESPONSE_OK)
		return
	}

	var user User
	has, err := GetUserFromDbByID(&user)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, RESPONSE_OK)
		return
	}

	if !has {
		WriteResponseWithCode(c, "账号不存在", nil, RESPONSE_OK)
		return
	}

	if user.Passwd != GenMD5WithSalt(password, user.Salt) {
		WriteResponseWithCode(c, "密码不正确", nil, RESPONSE_OK)
		return
	}

	//账密验证通过，生成session
	session := NewSession(user)
	err := session.Store() // 存储session到redis
	if err != nil {
		WriteResponseWithCode(c, "登录失败，请重试", nil, RESPONSE_OK)
		return
	}

	//登录成功，重定向到首页
	c.SetCookie("SESSION", session.ID, 0, "", "", false, true)
	c.SetCookie("USERNAME", session.Username, 0, "", "", false, true)
	c.SetCookie("UID", session.UID, 0, "", "", false, true)
	c.Redirect(http.StatusFound, MAIN_WEBSITE) // 重定向跳回首页

}

func GetCtxUser(c *gin.Context) *User {
	return c.Value("user").(*User)
}

func SignOut(c *gin.Context) {
	user := GetCtxUser(c)
	if user == nil {
		WriteResponseWithCode(c, "尚未登录", err.Error(), RESPONSE_OK)
		return
	}

	err = user.Del()
	if err != nil {
		WriteResponseWithCode(c, "注销失败，请重试", err.Error(), RESPONSE_OK)
		return
	}

	WriteResponseWithCode(c, "注销成功", nil, RESPONSE_OK)
	c.SetCookie("SESSION", "", 0, "", "", false, true)
	c.SetCookie("USERNAME", "", 0, "", "", false, true)
	c.SetCookie("UID", "", 0, "", "", false, true)

	c.Redirect(http.StatusFound, MAIN_WEBSITE)
}

func IsValidName(username string) error {
	return nil
}

func IsValidEmail(email string) error {
	return nil
}

func IsValidPasswd(passwd string) error {
	return nil
}


// 注册
func SignUp(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	retryPassword := c.PostForm("retry_password")
	email := c.PostForm("email")

	if retryPassword != password {
		WriteResponseWithCode(c, "两次输入的密码不一致", nil, RESPONSE_OK)
		return
	}

	err := IsValidName(username)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, RESPONSE_OK)
		return
	}

	err = IsValidEmail(email)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, RESPONSE_OK)
		return
	}

	err = IsValidPasswd(password)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, RESPONSE_OK)
		return
	}


	if UserExists("username", username) {
		WriteResponseWithCode(c, "用户名已存在", nil, RESPONSE_OK)
		return
	}
	if UserExists("email", email) {
		WriteResponseWithCode(c, "邮箱已被注册", nil, RESPONSE_OK)
		return
	}


	user := CreateUser(username, password, email)

	// 插入数据库
	err := DbInsertNewUser(user)
	if err != nil {
		WriteResponseWithCode(c, "注册失败，请稍后重新注册", nil, RESPONSE_OK)
		return
	}
	//注册完成后是自动登录跳转首页还是跳转到登录页面
	if SignInAfterSignUp {
		SetLoginCookie(ctx, username)
		WriteResponseWithCode(c, "", rsp, RESPONSE_OK)
	} else {
		WriteResponseWithCode(c, "", rsp, RESPONSE_OK)
	}
}


// 忘记密码
func ForgetPasswd(c *gin.Context) {


}

// 更新密码
func UpdatePasswd(c *gin.Context) {
	user := GetCtxUser(c)
	if user == nil {
		WriteResponseWithCode(c, "尚未登录", err.Error(), RESPONSE_OK)
		return
	}

	password := c.PostForm("password")

	err = IsValidPasswd(password)
	if err != nil {
		WriteResponseWithCode(c, err.Error(), nil, RESPONSE_OK)
		return
	}

	var updates = map[string]interface{} {
		"passwd":password
	}

	err := DbUpdateUser(updates)
	if err != nil {
		WriteResponseWithCode(c, "修改密码失败，请重试", nil, RESPONSE_OK)
		return
	}

	WriteResponseWithCode(c, "修改密码成功", nil, RESPONSE_OK)
}