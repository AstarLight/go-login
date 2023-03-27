package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

// 鉴权，登录态管理

var (
	ctx     = context.Background()
	RedisDb *redis.Client
)

func RedisInit() {
	RedisDb = redis.NewClient(&redis.Options{
		Addr: Conf.Redis.Addr,
	})

	pong, err := RedisDb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("connect redis fail: %v", err))
	} else {
		fmt.Println("connect redis succ,", pong)
	}

}

// session 登录态管理
type Session struct {
	ID         string `json:"id"` // sessionid
	Username   string `json:"username"`
	CreateTime int64  `json:"create_time"`
	IsAdmin    bool   `json:"is_admin"`
	UID        int64  `json:"uid"`
}

func GetUserFromSession(c *gin.Context) (error, *Session) {
	sessionid, _ := c.Cookie("SESSION")
	username, _ := c.Cookie("USERNAME")
	if len(sessionid) <= 0 || len(username) <= 0 {
		return errors.New("cookie session not exist"), nil
	}
	sess := GetSession(username)
	if sess == nil {
		return errors.New("session not exist"), nil
	}

	return nil, sess
}

func NewSession(user *User) *Session {
	now := time.Now()
	var newSession Session
	newSession.Username = user.Name
	newSession.ID = uuid.New().String()
	newSession.CreateTime = now.Unix()
	newSession.UID = user.Uid
	newSession.IsAdmin = user.IsAdmin
	return &newSession
}

func GetSession(username string) *Session {
	key := GetSessionKey(username)
	val, err := RedisDb.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("redis get失败，错误信息：%v\n", err)
		return nil
	}

	var userSession Session
	json.Unmarshal([]byte(val), &userSession)

	return &userSession

}

func GetSessionKey(username string) string {
	return Conf.Session.Key + username
}

func (s *Session) IsAdminUser() bool {
	return s.IsAdmin
}

// session信息存储到redis
func (s *Session) Store() error {
	key := GetSessionKey(s.Username)
	jdata, _ := json.Marshal(s)
	err := RedisDb.Set(ctx, key, string(jdata), time.Duration(Conf.Session.TTL)*time.Second).Err()
	if err != nil {
		fmt.Printf("redis set fail %v\n", err)
		return err
	}

	fmt.Println("session set key ", key)
	fmt.Println("session set val ", string(jdata))

	return nil
}

// 删除Session
func (s *Session) Del() error {
	key := GetSessionKey(s.Username)
	err := RedisDb.Del(ctx, key).Err()
	if err != nil {
		fmt.Printf("redis del fail %v\n", err)
		return err
	}

	fmt.Println("session del key ", key)

	return nil
}

func DelSession(username string) error {
	key := GetSessionKey(username)
	err := RedisDb.Del(ctx, key).Err()
	if err != nil {
		fmt.Printf("redis del fail %v\n", err)
		return err
	}

	fmt.Println("DelSession key ", key)

	return nil
}
