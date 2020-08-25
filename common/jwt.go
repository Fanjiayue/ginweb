package common

import (
	"github.com/dgrijalva/jwt-go"
	"ginweb/model"
	"time"
)

var jwtKey = []byte("secret_create")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}


func ReleaseToken(user model.User)(string,error){
	exirationTime := time.Now().Add(7*24*time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt:exirationTime.Unix(),
			IssuedAt:time.Now().Unix(),
			Issuer:"oceanlearn.tech",
			Subject:"user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "",err
	}

	return tokenString,nil
}

//解析token
func ParseToken(tokenString string) (*jwt.Token,*Claims,error){
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token,claims,err
}