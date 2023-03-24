package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var MasterDB *xorm.Engine


var (
	ConnectDBErr = errors.New("connect db error")
	UseDBErr     = errors.New("use db error")

	WriteDBErr     = errors.New("write but affect zero row")
)


func init() {
	mysqlConfig, err := ConfigFile.GetSection("mysql")
	if err != nil {
		fmt.Println("get mysql config error:", err)
		return
	}


	// 启动时就打开数据库连接
	if err = initEngine(); err != nil {
		panic(err)
	}

	// 测试数据库连接是否 OK
	if err = MasterDB.Ping(); err != nil {
		panic(err)
	}
}

func initEngine() error {
	var err error

	MasterDB, err = xorm.NewEngine("mysql", "lijunshi:lijunshipwd@tcp(10.10.40.231:3306)/users")
	if err != nil {
		return err
	}

	maxIdle := ConfigFile.MustInt("mysql", "max_idle", 2)
	maxConn := ConfigFile.MustInt("mysql", "max_conn", 10)

	MasterDB.SetMaxIdleConns(maxIdle)
	MasterDB.SetMaxOpenConns(maxConn)

	showSQL := ConfigFile.MustBool("xorm", "show_sql", false)
	logLevel := ConfigFile.MustInt("xorm", "log_level", 1)

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

func (InstallLogic) CreateTable() error {

	dbFile := config.ROOT + "/config/db.sql"
	buf, err := ioutil.ReadFile(dbFile)

	if err != nil {
		objLog.Errorln("create table, read db file error:", err)
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
		}
	}

	return err
}

func GetUserFromDbByName(user *User) (bool, error) {
	has, err := MasterDB.Where("Name=?", user.Name).Get(user)
	if err != nil {
		// 发生错误，认为已经创建了
		return false, err
	}

	return has, nil
}

func GetUserFromDbByID(user *User) (bool, error) {
	has, err := MasterDB.Where("ID=?", user.ID).Get(user)
	if err != nil {
		// 发生错误，认为已经创建了
		return false, err
	}

	return has, nil
}

func DbInsertNewUser(user *User) (error) {
	affected, err := engine.Insert(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return WriteDBErr
	}
	return nil
}

func DbUpdateUser(updates map[string]interface{}) (error) {
	affected, err := engine.Table("user").Update(updates)
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

	user := &User{}
	has, err := MasterDB.Where(field+"=?", val).Get(user)
	if err != nil || user.ID == 0 || !has {
		if err != nil {
			objLog.Errorln("user logic UserExists error:", err)
		}
		return false
	}
	return true
}