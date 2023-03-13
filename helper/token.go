package helper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// 用户声明，用来创建用户token
type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

// 用来保存token的key
var myKey = []byte("gin-gorm-oj-key")

// GenerateToken
// 根据用户唯一标识和用户名，生成token
func GenerateToken(identity string, name string, isAdmin int) (string, error) {

	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim) // 封装的token对象
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// AnalyseToken
// 解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error: %v", err)
	}
	return userClaim, nil
}
