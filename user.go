package main

import (
	"fmt"
	"math/rand"
	"time"
)

type User struct {
	ID     int64  `xorm:"pk autoincr"`     //主键，UID
	Name   string `xorm:"UNIQUE NOT NULL"` // 用户名，唯一键
	Email  string `xorm:"UNIQUE NOT NULL"` // 邮箱，唯一键
	Passwd string `xorm:"NOT NULL"`        //已使用Salt进行加密的密码串MD5(原始password+Salt)

	Salt string `xorm:"VARCHAR(32)"` // 用于密码加盐哈希，注册时随机生成

	CreatedUnix   int64 `xorm:"INDEX created"` // 账号创建时间
	UpdatedUnix   int64 `xorm:"INDEX updated"` // 账号更新时间
	LastLoginUnix int64 `xorm:"INDEX"`         // 账号上次登录时间

	IsAdmin bool `xorm:"NOT NULL DEFAULT false"` // 管理员标记

	ProhibitLogin bool `xorm:"NOT NULL DEFAULT false"` // 禁止登录标记

	LastLoginIp string `xorm:"VARCHAR(32) INDEX"`
}

func CreateUser(username, password, email string) *User {
	newUser := User{}
	newUser.Email = email
	newUser.Name = username
	newUser.Salt = fmt.Sprintf("%x", rand.Int31())
	newUser.Passwd = GenMD5WithSalt(password, newUser.Salt)
	newUser.CreatedUnix = time.Now().Unix()
	newUser.UpdatedUnix = time.Now().Unix()
	newUser.IsAdmin = false
	newUser.ProhibitLogin = false
	return &newUser
}

func GetUser(user *User) (bool, error) {
	userInLocal, err := GetUserFromLocal(user)
	if err != nil {
		return false, err
	}

	if !userInLocal {
		userInRedis, err := GetUserFromRedis(user)
		if err != nil {
			return false, err
		}

		if !userInRedis {
			userInDB, err := GetUserFromDb(user)
			if err != nil {
				return false, err
			}

			if !userInDB {
				return false, nil
			} else {
				SetUserToRedis(user)
				SetUserToLocal(user)
			}

		} else {
			SetUserToLocal(user)
		}
	}

	return true, nil

}

func GetUserFromDb(user *User) (bool, error) {
	return false, nil
}

func GetUserFromRedis(user *User) (bool, error) {
	return false, nil
}

func GetUserFromLocal(user *User) (bool, error) {
	return false, nil
}

func SetUserToRedis(user *User) error {
	return nil
}

func SetUserToLocal(user *User) error {
	return nil
}
