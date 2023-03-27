package main

import (
	//"database/sql"
	"bytes"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var MasterDB *xorm.Engine

var (
	ConnectDBErr = errors.New("connect db error")
	UseDBErr     = errors.New("use db error")

	WriteDBErr = errors.New("write but affect no row")
)

func DbInit() {
	// 启动时就打开数据库连接
	if err := initEngine(); err != nil {
		panic(err)
	}

	// 测试数据库连接是否 OK
	if err := MasterDB.Ping(); err != nil {
		panic(err)
	}

	// user表不存在，则先建表
	if !IsTableExist() {
		err := CreateTable()
		if err != nil {
			panic(err)
		}
	}
}

func initEngine() error {
	var err error

	MasterDB, err = xorm.NewEngine("mysql", Conf.MySQL.Addr)
	if err != nil {
		return err
	}

	maxIdle := Conf.MySQL.MaxIdle
	maxConn := Conf.MySQL.MaxConn

	MasterDB.SetMaxIdleConns(maxIdle)
	MasterDB.SetMaxOpenConns(maxConn)

	showSQL := Conf.MySQL.ShowSQL
	logLevel := Conf.MySQL.LogLevel

	MasterDB.ShowSQL(showSQL)
	MasterDB.Logger().SetLevel(log.LogLevel(logLevel))

	// 启用缓存
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// MasterDB.SetDefaultCacher(cacher)

	return nil
}

func IsTableExist() bool {
	exists, err := MasterDB.IsTableExist(new(User))
	if err != nil || !exists {
		return false
	}

	return true
}

func CreateTable() error {

	dbFile := Conf.SQL.File
	buf, err := ioutil.ReadFile(dbFile)

	if err != nil {
		fmt.Println("create table, read db file error:", err)
		return err
	}

	sqlSlice := bytes.Split(buf, []byte("CREATE TABLE"))
	MasterDB.Exec("SET SQL_MODE='ALLOW_INVALID_DATES';")
	for _, oneSql := range sqlSlice {
		strSql := string(bytes.TrimSpace(oneSql))
		if strSql == "" {
			continue
		}

		strSql = "CREATE TABLE " + strSql
		_, err1 := MasterDB.Exec(strSql)
		if err1 != nil {
			fmt.Println("create table error:", err1)
			err = err1
		} else {
			fmt.Println("create table succ, sql=", strSql)
		}
	}

	return err
}

func GetUserFromDbByName(user *User) (bool, error) {
	has, err := MasterDB.Where("name=?", user.Name).Get(user)
	if err != nil {
		return false, err
	}

	return has, nil
}

// 向数据库插入一个新用户数据
func DBInsertNewUser(user *User) error {
	affected, err := MasterDB.Insert(user)
	if err != nil {
		fmt.Println("DBInsertNewUser err ", err)
		return err
	}
	if affected == 0 {
		return WriteDBErr
	}
	return nil
}

func DBUpdateUser(username string, updates map[string]interface{}) error {
	updates["updated_unix"] = time.Now().Unix()
	affected, err := MasterDB.Table("user").Where("name=?", username).Update(updates)
	if err != nil {
		return err
	}
	if affected == 0 {
		return WriteDBErr
	}
	return nil

}

// UserExists 判断用户是否存在
func UserExists(field, val string) bool {
	fmt.Println(field, val)
	user := &User{}
	has, err := MasterDB.Where(field+"=?", val).Get(user)
	if err != nil || user.Uid == 0 {
		if err != nil {
			fmt.Println("user logic UserExists error:", err)
		}
		return false
	}
	return has
}
