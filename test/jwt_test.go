package test

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type UserClaims struct {
	Identity string
	Name     string
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// 生成token
func TestGenerateToken(t *testing.T) {
	UserClaim := &UserClaims{
		Identity:       "user_1",
		Name:           "Get",
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim) // 封装的token对象
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenString)

}

// 解析token
func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZGVudGl0eSI6InVzZXJfMSIsIk5hbWUiOiJHZXQifQ.O8q7NCOLZp9Bgk-qoPQ68eE5N7r5jzEZ1tvmXVly7u4"
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(userClaim)
	}
}

// 通过MD5加密
func TestGenerateMd5(t *testing.T) {
	var password string = "bcs@1234"

	psw := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	fmt.Println(psw)
	var err error = nil
	if err != nil {
		t.Fatal(err)
	}
}
