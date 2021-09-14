package common

import (
	"github.com/dgrijalva/jwt-go"
	"oceanlearn.learn/ginessential/model"
	"time"
)

var jwtKey = []byte("a_secret_kiki")

type Claims struct{
	UserId uint
	jwt.StandardClaims
}

//发放token
func ReleaseToken(user model.User) (string, error) {
	//token发放持续时间，是说免登录保持时长好像
	expirationTime := time.Now().Add(7*24*time.Hour)
	claims := &Claims{
		UserId : user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "flare",
			Subject: "user token",
		},

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString ,err := token.SignedString(jwtKey)
	if err!=nil{
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error){
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token)(i interface{}, err error){
		return jwtKey, nil
	})
	return token, Claims, err
}