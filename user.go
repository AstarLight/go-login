package main

import (
	"time"
)

type User struct {
	Uid    int64  `xorm:"pk autoincr"`     //主键，UID
	Name   string `xorm:"UNIQUE NOT NULL"` // 用户名，唯一键
	Email  string `xorm:"UNIQUE NOT NULL"` // 邮箱，唯一键
	Passwd string `xorm:"NOT NULL"`        //已使用Salt进行加密的密码串MD5(原始password+Salt)

	Salt string `xorm:"VARCHAR(32)"` // 用于密码加盐哈希，注册时随机生成

	CreatedUnix   int64 `xorm:"INDEX created"` // 账号创建时间
	UpdatedUnix   int64 `xorm:"INDEX updated"` // 账号更新时间
	LastLoginUnix int64 `xorm:"INDEX"`         // 账号上次登录时间

	IsAdmin bool `xorm:"NOT NULL DEFAULT false"` // 管理员标记

	ProhibitLogin bool `xorm:"NOT NULL DEFAULT false"` // 禁止登录标记

	LastLoginIp string `xorm:"VARCHAR(32) INDEX"` // 登录Ip
}

func (this *User) TableName() string {
	return "user"
}

func CreateUser(username, password, email string) *User {
	newUser := User{}
	newUser.Email = email
	newUser.Name = username
	newUser.Salt = GenRandomSalt()
	newUser.Passwd = GenMD5WithSalt(password, newUser.Salt)
	newUser.CreatedUnix = time.Now().Unix()
	newUser.UpdatedUnix = time.Now().Unix()
	newUser.IsAdmin = false
	newUser.ProhibitLogin = false
	return &newUser
}
