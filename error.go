package yljwt

const (
	ErrBadPosition JWTError = iota + 1
	ErrNotFindCookie
	ErrTokenOut
	ErrNoUserID
	ErrAppIDNotAllow
)

//JWTError 错误
type JWTError int

var errMap = map[JWTError]string{
	ErrBadPosition:   "使用了未定义的token位置",
	ErrNotFindCookie: "cookie中没有找到token，检查cookie中是否包含token节点",
	ErrTokenOut:      "token已经过期",
	ErrNoUserID:      "token中没有找到UserID",
	ErrAppIDNotAllow："请求的来源不允许",
}

func (j JWTError) Error() string {
	str := errMap[j]
	if str == "" {
		str = "未定义的错误类型"
	}
	return str
}
