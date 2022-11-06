package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已过期",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "无Token数据",
	ERROR_USER_EXISTS:              "用户已存在",
	ERROR_USER_CREATE:              "用户注册失败",
	ERROR_USER_PASSWORD:            "用户名或密码错误",
}

func GetMsg(errCode int) string {
	msg, ok := MsgFlags[errCode]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
