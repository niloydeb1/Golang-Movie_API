package v1

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/niloydeb1/Golang-Movie_API/api/common"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"log"
	"time"
)

type userApi struct {
}

// Registration... Registration Api
// @Summary Registration api
// @Description Api for users registration
// @Tags User
// @Produce json
// @Param Authorization header string true "Insert your access token while adding new user" default(Bearer <Add access token here>)
// @Param data body v1.UserRegistrationDto true "dto for creating user"
// @Param action path string false "action [create_admin] if superadmin wants to create new admin"
// @Success 200 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/users [POST]
func (u userApi) Registration(context echo.Context) error {
	registrationType := context.QueryParam("action")
	if registrationType == "" {
		return u.registerAdmin(context)
	} else if registrationType == string(enums.CREATE_ADMIN) {
		return u.registerUser(context)
	}
	return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", errors.New("invalid query action").Error())
}

func (u userApi) registerAdmin(context echo.Context) error {
	userFromToken, err := GetUserFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.Role != enums.SUPERADMIN {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	formData := v1.UserRegistrationDto{}
	if err = context.Bind(&formData); err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, nil, "Failed to Bind Input!")
	}
	if formData.Password == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", "password is required")
	} else if len(formData.Password) < 8 {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", "password length must be at least 8")
	}
	formData.ID = uuid.New().String()
	formData.CreatedDate = time.Now().UTC()
	formData.UpdatedDate = time.Now().UTC()
	formData.Status = enums.ACTIVE
	err = formData.Validate()
	if err != nil {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", err.Error())
	}
	user :=v1.GetUserFromUserRegistrationDto(formData)
	user.Role = enums.ADMIN
	userExist := v1.User{}.GetByEmail(user.Email)
	if userExist.Email != "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", "Email is already registered.")
	}
	err = v1.User{}.Store(user)
	if err != nil {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", err.Error())
	}
	return common.GenerateSuccessResponse(context, formData, nil, "Successfully Created User!")
}

func (u userApi) registerUser(context echo.Context) error {
	formData := v1.UserRegistrationDto{}
	if err := context.Bind(&formData); err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, nil, "Failed to Bind Input!")
	}
	if formData.Password == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", "password is required")
	} else if len(formData.Password) < 8 {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", "password length must be at least 8")
	}
	formData.ID = uuid.New().String()
	formData.CreatedDate = time.Now().UTC()
	formData.UpdatedDate = time.Now().UTC()
	formData.Status = enums.ACTIVE
	err := formData.Validate()
	if err != nil {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", err.Error())
	}
	user :=v1.GetUserFromUserRegistrationDto(formData)
	user.Role = enums.USER
	userExist := v1.User{}.GetByEmail(user.Email)
	if userExist.Email != "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", "Email is already registered.")
	}
	err = v1.User{}.Store(user)
	if err != nil {
		return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", err.Error())
	}
	return common.GenerateSuccessResponse(context, formData, nil, "Successfully Created User!")
}
