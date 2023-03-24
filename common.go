package main



func GenMD5WithSalt(passwd, salt string) string {
	s := passwd + "::" + salt
	md5Hash := md5.New()
	md5Hash.Write([]byte(s))
	// 转16进制
	return hex.EncodeToString(md5Hash.Sum(nil))
}


type Response struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Data interface{} `json:"data"`

	// Reference returns the reference document which maybe useful to solve this error.
	Msg string `json:"msg"`
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponseWithCode(c *gin.Context, msg string, data interface{}, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// 时间戳转时间 2022-02-01 12:00:00
func TimestampParse(t int64) string {
	tm := time.Unix(t, 0)
	return tm.Format("2006-01-02 15:04:05")
}

// 时间转时间戳 2023-03-23 15:55:00
func DateToTimestamp(t string) int64 {
	loc, _ := time.LoadLocation("Local") //获取当地时区
	ts, err := time.ParseInLocation("2006-01-02 15:04:05", t, loc)
	if err != nil {
	   return 0
	}

	return ts.Unix()
}