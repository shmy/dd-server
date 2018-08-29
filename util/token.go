package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

func GenerateTheToken(id interface{}, platform string) (string, error) {
	var overdue = viper.GetInt64("jsonwebtoken." + platform + ".overdue")
	var secret = viper.GetString("jsonwebtoken." + platform + ".secret")
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Second * time.Duration(overdue)).Unix()
	return token.SignedString([]byte(secret))
}
