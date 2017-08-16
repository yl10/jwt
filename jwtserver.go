package yljwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

//var spaceuuid=uuid.NewMD5()
const (
	Header TokenPosition = iota + 1
	Cookie
	Query
)

//TokenPosition token位置
type TokenPosition int

//JwtServerGroup server组
type JwtServerGroup map[string]*JwtServer

//JwtServer server
type JwtServer struct {
	appID    string
	key      interface{}
	position TokenPosition
	method   jwt.SigningMethod
	expires  time.Duration
}

//NewJwtServer 创建一个新JWT服务
func NewJwtServer(appid string, key interface{}, position TokenPosition, method jwt.SigningMethod, expires time.Duration) *JwtServer {
	return &JwtServer{
		appid, key, position, method, expires,
	}
}

//NewJwtServerHS256 创建一个新HS256JWT服务
func NewJwtServerHS256(appid string, key string, position TokenPosition, expires time.Duration) *JwtServer {
	return &JwtServer{
		appid, []byte(key), position, jwt.SigningMethodHS256, expires,
	}
}

//UserInfo 用户信息
type UserInfo struct {
	LoginKind string `json:"loginkind"`
	LoginID   string `json:"loginid"`
	LoginFrom string `json:"loginfrom"`
}

func (u UserInfo) GetKind() string {
	return u.LoginKind
}

func (u UserInfo) GetID() string {
	return u.LoginID
}

//LoginUser 登录用户接口
type LoginUser interface {
	GetKind() string
	GetID() string
	// CheckUser() bool
}

//YlClaims yl声明
type YlClaims struct {
	*jwt.StandardClaims
	UserInfo
}

//GetID 获取ID
func (y YlClaims) GetID() string {
	return y.Id
}

//NewYlClaims 创建一个声明
func (j JwtServer) newYlClaims(user LoginUser) YlClaims {

	c := YlClaims{

		&jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + int64(j.expires),
			Id:        uuid.New(),
		},
		UserInfo{
			LoginKind: user.GetKind(),
			LoginID:   user.GetID(),
			LoginFrom: j.appID,
		},
	}

	return c
}

//Token 生成一个token字符串
func (j JwtServer) Token(user LoginUser, store TokenStore) (tokenstring string, err error) {
	c := j.newYlClaims(user)
	token := jwt.NewWithClaims(j.method, c)

	tokenstring, err = token.SignedString(j.key)
	if store != nil {
		err = store.Create(c.Id)
	}
	return
}

//WriteToken 直接输出token,放在body中，｛token:"token值"｝,如果是cooke，直接防在token节点中
func (j JwtServer) WriteToken(w http.ResponseWriter, user LoginUser, store TokenStore) error {

	str, err := j.Token(user, store)
	if err != nil {
		return err
	}

	switch j.position {
	case Cookie:
		ck := &http.Cookie{Name: "token", Value: str, Expires: time.Now().Add(j.expires)}
		http.SetCookie(w, ck)
		return nil
	default:
		_, err = w.Write([]byte(str))
		return err
	}

}

//CheckToken 通过http请求，检查token
func (j JwtServer) CheckToken(r *http.Request, checkuser func(user UserInfo) bool) (user UserInfo, err error) {
	var tokenstring string
	switch j.position {
	case Header:
		tokenstring = r.Header.Get("Authorization")
	case Query:
		r.ParseForm()
		tokenstring = r.PostFormValue("token")
	case Cookie:
		ck, err2 := r.Cookie("token")
		if err2 != nil {
			err = ErrNotFindCookie
			return
		}
		tokenstring = ck.Value

	default:
		err = ErrBadPosition
		return
	}

	// token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
	// 	return j.key, nil
	// })

	token, err := jwt.ParseWithClaims(tokenstring, &YlClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		return
	}
	yc := token.Claims.(*YlClaims)
	//检查过期
	if yc.VerifyExpiresAt(time.Now().Unix(), true) == false {
		err = ErrTokenOut
		return
	}

	//检查appid
	if yc.LoginFrom != j.appID {
		err = ErrAppIDNotAllow
		return
	}
	user.LoginFrom = yc.LoginFrom
	user.LoginID = yc.LoginID
	user.LoginKind = yc.LoginKind

	if checkuser != nil {
		if !checkuser(user) {
			err = fmt.Errorf("用户校验未通过")
			return
		}
	}
	return

}

// //Claims 声明 包含标准声明和扩展部分
// type Claims struct {
// 	StandardClaims jwt.StandardClaims
// 	ExtendedClaims map[string]interface{}
// }

// //SetSubject 设置主题
// func (c *Claims) SetSubject(subject string) *Claims {
// 	c.StandardClaims.Subject = subject
// 	return c
// }

// //SetAudience 设置读者
// func (c *Claims) SetAudience(audience string) *Claims {
// 	c.StandardClaims.Audience = audience
// 	return c
// }

// //SetIssuer 设置发行者
// func (c *Claims) SetIssuer(issuer string) *Claims {
// 	c.StandardClaims.Issuer = issuer
// 	return c
// }

// //SetNotBefore 设置生效时间
// func (c *Claims) SetNotBefore(notbefore int64) *Claims {
// 	c.StandardClaims.NotBefore = notbefore
// 	return c
// }

// //SetExpiresAt 设置失效时间
// func (c *Claims) SetExpiresAt(expiresAt int64) *Claims {
// 	c.StandardClaims.ExpiresAt = expiresAt
// 	return c
// }

// //SetExtend 设置自定义声明
// func (c *Claims) SetExtend(k string, v interface{}) *Claims {
// 	if c.ExtendedClaims == nil {
// 		c.ExtendedClaims = make(map[string]interface{})
// 	}
// 	c.ExtendedClaims[k] = v
// 	return c
// }

// func (d DefaultTokenStore)
// //生成token
// func (j JwtServer)Token(){
// token:=jwt.New(j.method)
// token.
// }
