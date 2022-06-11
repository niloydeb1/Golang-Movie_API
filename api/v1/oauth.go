package v1

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/niloydeb1/Golang-Movie_API/api/common"
	"github.com/niloydeb1/Golang-Movie_API/config"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

// OauthRouter api/v1/oauth/* router
func OauthRouter(g *echo.Group) {
	g.POST("/login", oauthApi{}.Login)
}

type oauthApi struct {
}

// Login... Login Api
// @Summary Login api
// @Description Api for users login
// @Tags Oauth
// @Produce json
// @Param loginData body v1.LoginDto true "Login dto if grant_type=password"
// @Param refreshTokenData body v1.RefreshTokenDto true "RefreshTokenDto dto if grant_type=refresh_token"
// @Success 200 {object} common.ResponseDTO{data=v1.JWTPayLoad{}}
// @Failure 403 {object} common.ResponseDTO
// @Router /api/v1/oauth/login [POST]
func (o oauthApi) Login(context echo.Context) error {
	if context.QueryParam("grant_type") == "password" {
		return o.handlePasswordGrant(context)
	} else if context.QueryParam("grant_type") == "refresh_token" {
		return o.handleRefreshTokenGrant(context)
	}
	return common.GenerateForbiddenResponse(context, nil, "Please provide a valid grant_type")
}

func (o oauthApi) handleRefreshTokenGrant(context echo.Context) error {
	refreshTokenDto := new(v1.RefreshTokenDto)
	if err := context.Bind(&refreshTokenDto); err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, "[ERROR]: Failed bind payload from context", err.Error())
	}
	tokenValid := v1.Jwt{}.IsTokenValid(refreshTokenDto.RefreshToken)
	if !tokenValid {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Token is expired!", "Please login again to get token!")
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(refreshTokenDto.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Publickey), nil
	})
	jsonbody, err := json.Marshal(claims["data"])
	if err != nil {
		log.Println(err)
	}
	userFromToken := v1.UserTokenDto{}
	if err := json.Unmarshal(jsonbody, &userFromToken); err != nil {
		log.Println(err)
	}
	existingUser := v1.User{}.GetByID(userFromToken.ID)
	if existingUser.ID == "" || existingUser.Status != enums.ACTIVE {
		return common.GenerateForbiddenResponse(context, "[ERROR]: No User found!", "Please login with actual user email!")
	}
	tokenLifeTime, err := strconv.ParseInt(config.TokenLifetime, 10, 64)
	if err != nil {
		log.Println(err.Error())
		return common.GenerateForbiddenResponse(context, "[ERROR]: failed to read regular token lifetime from env!", err.Error())
	}
	token, refreshToken, err := v1.Jwt{}.GenerateToken(userFromToken.ID, tokenLifeTime, userFromToken)
	if err != nil {
		log.Println(err.Error())
		return common.GenerateForbiddenResponse(context, "[ERROR]: failed to create token!", err.Error())
	}

	err = v1.TokenService{}.Store(v1.Token{Uid: userFromToken.ID, Token: token, RefreshToken: refreshToken})
	if err != nil {
		log.Println(err.Error())
		return common.GenerateForbiddenResponse(context, "[ERROR]: failed to store token!", err.Error())
	}
	return common.GenerateSuccessResponse(context, v1.JWTPayLoad{AccessToken: token, RefreshToken: refreshToken}, nil, "")
}

func (o oauthApi) handlePasswordGrant(context echo.Context) error {
	loginDto := new(v1.LoginDto)
	if err := context.Bind(&loginDto); err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, "[ERROR]: Failed bind payload from context", err.Error())
	}

	existingUser := v1.User{}.GetByEmail(loginDto.Email)
	if existingUser.ID == "" || existingUser.Status != enums.ACTIVE {
		return common.GenerateForbiddenResponse(context, "[ERROR]: No User found!", "Please login with actual user email!")
	}
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginDto.Password))
	if err != nil {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Password not matched!", "Please login with due credential!"+err.Error())
	}
	tokenLifeTime, err := strconv.ParseInt(config.TokenLifetime, 10, 64)
	if err != nil {
		log.Println(err.Error())
		return common.GenerateForbiddenResponse(context, "[ERROR]: failed to read regular token lifetime from env!", err.Error())
	}
	userTokenDto := v1.UserTokenDto{
		ID:        existingUser.ID,
		FirstName: existingUser.FirstName,
		LastName:  existingUser.LastName,
		Email:     existingUser.Email,
		Phone:     existingUser.Phone,
		Status:    existingUser.Status,
		Role:      existingUser.Role,
	}
	token, refreshToken, err := v1.Jwt{}.GenerateToken(userTokenDto.ID, tokenLifeTime, userTokenDto)
	if err != nil {
		log.Println(err.Error())
		return common.GenerateForbiddenResponse(context, "[ERROR]: failed to create token!", err.Error())
	}

	err = v1.TokenService{}.Store(v1.Token{Uid: userTokenDto.ID, Token: token, RefreshToken: refreshToken})
	if err != nil {
		log.Println(err.Error())
		return common.GenerateForbiddenResponse(context, "[ERROR]: failed to store token!", err.Error())
	}
	return common.GenerateSuccessResponse(context, v1.JWTPayLoad{AccessToken: token, RefreshToken: refreshToken}, nil, "")
}