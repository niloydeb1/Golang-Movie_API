package v1

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/niloydeb1/Golang-Movie_API/config"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"log"
	"strings"
)

// GetUserFromBearerToken returns user from bearer token
func GetUserTokenDtoFromBearerToken(context echo.Context, jwtService v1.Jwt) (v1.UserTokenDto, error) {
	bearerToken := context.Request().Header.Get("Authorization")
	if bearerToken == "" {
		return v1.UserTokenDto{}, errors.New("[ERROR]: No token found!")
	}
	var token string
	if len(strings.Split(bearerToken, " ")) == 2 {
		token = strings.Split(bearerToken, " ")[1]
	} else {
		return v1.UserTokenDto{}, errors.New("[ERROR]: No token found!")
	}
	if !jwtService.IsTokenValid(token) {
		return v1.UserTokenDto{}, errors.New("[ERROR]: Token is expired!")
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Publickey), nil
	})
	jsonbody, err := json.Marshal(claims["data"])
	if err != nil {
		log.Println(err)
	}
	userTokenDto := v1.UserTokenDto{}
	if err := json.Unmarshal(jsonbody, &userTokenDto); err != nil {
		return v1.UserTokenDto{}, errors.New("[ERROR]: No User Found!")
	}
	return userTokenDto, nil
}