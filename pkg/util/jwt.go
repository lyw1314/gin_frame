package util

import (
	"gin_frame/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	jwtSecret []byte
	Issuer    string
	nowFunc   func() time.Time
}

//	==Payload默认7个字段==
//	Audience  接收方
//	ExpiresAt jwt过期时间
//	Id        jwt唯一身份标识
//	IssuedAt  签发时间
//	Issuer    签发人/发行人
//	NotBefore jwt生效时间
//	Subject   主题

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewJwt(issuer string) *JWT {
	return &JWT{
		jwtSecret: []byte(setting.AppC.GetString("app.JWT_SECRET")),
		Issuer:    issuer,
		nowFunc:   time.Now,
	}
}

// CreateToken 生成token
func (j *JWT) CreateToken(userName string, accessID string, expire time.Duration) (string, error) {
	nowTime := j.nowFunc().Unix()
	var expireTime int64
	if expire > 0 {
		expireTime = nowTime + int64(expire.Seconds())
	}
	withClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			IssuedAt:  nowTime,
			Issuer:    j.Issuer,
			Subject:   accessID,
		},
	})
	signedString, err := withClaims.SignedString(j.jwtSecret)
	return signedString, err
}

// ParseToken 解析token
func (j *JWT) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
