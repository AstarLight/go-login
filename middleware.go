package main

//必须登录的请求，从session读user写入context 
func NeedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, user := GetUserFromSession(c)
		if err == nil && user != nil {
			context.WithValue(c, "user", user)
			user.Store() // 已登录，每次请求都会续期
			c.Next()
			return

		} else {
			// 未登录
			WriteResponseWithCode(c, "未登录", nil, 429)
			c.Abort()
			return
		}
	}
}

//限流，粒度分为IP,MAC-ADDR
func CommonRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//黑名单,限流，粒度分为IP,MAC-ADDR
func CommonBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//限流，粒度分为UID
func UIDRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//黑名单,限流，粒度分为UID
func UIDBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}