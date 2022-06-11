package v1

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/niloydeb1/Golang-Movie_API/api/common"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func UserRouter(g *echo.Group) {
	g.POST("", userApi{}.Registration)
	g.GET("", userApi{}.Get)
	g.GET("/:id", userApi{}.GetByID)
	g.DELETE("/:id", userApi{}.Delete)
	g.PUT("", userApi{}.Update)
}

type userApi struct {
}

// Get... Get Api
// @Summary Get api
// @Description Api for getiing all user by admin
// @Tags User
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param status path string true "status type [active/inactive]"
// @Success 200 {object} common.ResponseDTO{data=[]v1.User{}}
// @Forbidden 403 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Router /api/v1/users [GET]
func (u userApi) Get(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.Role != enums.SUPERADMIN && userFromToken.Role != enums.ADMIN {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	status := context.QueryParam("status")
	if status == string(enums.ACTIVE) {
		return common.GenerateSuccessResponse(context, v1.User{}.GetUsers(enums.STATUS(status)), nil, "Success!")
	} else if status == string(enums.INACTIVE) {
		return common.GenerateSuccessResponse(context, v1.User{}.GetUsers(enums.STATUS(status)), nil, "Success!")
	}
	return common.GenerateForbiddenResponse(context, "[ERROR]: No valid status found!", "Please provide a valid status.")
}

// GetByID... GetByID Api
// @Summary Registration api
// @Description Api for getiing user by id
// @Tags User
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path string true "id user id"
// @Success 200 {object} common.ResponseDTO{data=v1.User{}}
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/users/{id} [GET]
func (u userApi) GetByID(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	id := context.Param("id")
	if userFromToken.ID != id {
		if userFromToken.Role != enums.SUPERADMIN && userFromToken.Role != enums.ADMIN {
			return common.GenerateForbiddenResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
		}
	}
	data := v1.User{}.GetByID(id)
	if data.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: User Not Found!", "Please give a valid user id!")
	}
	return common.GenerateSuccessResponse(context, data, nil, "Success!")
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
		return u.registerUser(context)
	} else if registrationType == string(enums.CREATE_ADMIN) {
		return u.registerAdmin(context)
	}
	return common.GenerateErrorResponse(context, "[ERROR]: Failed to register user!", errors.New("invalid query action").Error())
}

func (u userApi) registerAdmin(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
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

// Update... Update Api
// @Summary Update api
// @Description Api for updating users object
// @Tags User
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param action path string true "action type [reset_password/update_status]"
// @Param status path string false "status type [inactive/active] if action update_status"
// @Param id path string false "updating users id, if action update_status"
// @Param password_reset_dto body v1.PasswordResetDto true "dto for resetting users password"
// @Success 200 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/users [PUT]
func (u userApi) Update(context echo.Context) error {
	action := context.QueryParam("action")
	if action == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: No action type is provided!", "Please provide a action type!")
	} else if action == string(enums.RESET_PASSWORD) {
		return u.ResetPassword(context)
	} else if action == string(enums.UPDATE_STATUS) {
		return u.UpdateStatus(context)
	}
	return common.GenerateErrorResponse(context, "[ERROR]: Invalid type is provided!", "Please provide a valid action type!")
}

func (u userApi) UpdateStatus(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.Role != enums.SUPERADMIN && userFromToken.Role != enums.ADMIN {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	status := context.QueryParam("status")
	if enums.STATUS(status) != enums.ACTIVE && enums.STATUS(status) != enums.INACTIVE {
		return common.GenerateErrorResponse(context, "[ERROR]: Invalid update status!", "Please provide a valid update status!")
	}
	userId := context.QueryParam("id")
	user := v1.User{}.GetByID(userId)
	if userFromToken.Role == enums.ADMIN && (user.Role == enums.ADMIN || user.Role == enums.SUPERADMIN) {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	if user.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: User not found!", "Please provide a valid user id!")
	}
	if user.Status == enums.DELETED {
		return common.GenerateErrorResponse(context, "[ERROR]: User not found!", "Please provide a valid user id!")
	}
	err = v1.User{}.UpdateStatus(userId, enums.STATUS(status))
	if err != nil {
		return common.GenerateForbiddenResponse(context, err.Error(), "Operation Failed!")
	}
	return common.GenerateSuccessResponse(context, nil, nil, "Operation Successful!")
}

func (u userApi) ResetPassword(context echo.Context) error {
	formData := v1.PasswordResetDto{}
	if err := context.Bind(&formData); err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, nil, "Failed to Bind Input!")
	}
	if formData.CurrentPassword == "" {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Failed to reset password!", "Please provide required data!")
	}
	var user v1.User
	if formData.Email == "" {
		userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
		if err != nil {
			return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
		}
		user = v1.User{}.GetByID(userFromToken.ID)
	} else {
		user = v1.User{}.GetByEmail(formData.Email)
	}
	if user.ID == "" || user.Status != enums.ACTIVE {
		return common.GenerateForbiddenResponse(context, "[ERROR]: No User found!", "Please login with actual user email!")
	}
	if formData.CurrentPassword != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formData.CurrentPassword))
		if err != nil {
			return common.GenerateForbiddenResponse(context, "[ERROR]: Password not matched!", "Please provide due credential!"+err.Error())
		}
	}
	user.Password = formData.NewPassword
	err := v1.User{}.UpdatePassword(user)
	if err != nil {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Failed to reset password!", err.Error())
	}
	return common.GenerateSuccessResponse(context, nil, nil, "Operation Successful!")
}

// Delete... Delete Api
// @Summary Delete api
// @Description Api to delete user
// @Tags User
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param id path string true "id user id"
// @Success 200 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/users [DELETE]
func (u userApi) Delete(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.Role != enums.SUPERADMIN && userFromToken.Role != enums.ADMIN {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	id := context.Param("id")
	user := v1.User{}.GetByID(id)
	if userFromToken.Role == enums.ADMIN && (user.Role == enums.ADMIN || user.Role == enums.SUPERADMIN) {
		return common.GenerateForbiddenResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	if user.ID == "" || user.Status != enums.ACTIVE {
		return common.GenerateErrorResponse(context, "[ERROR]: User not found!", "Please provide a valid user id!")
	}
	err = v1.User{}.Delete(id)
	if err != nil {
		return common.GenerateErrorResponse(context, nil, "Failed to Delete User!")
	}
	return common.GenerateSuccessResponse(context, nil, nil, "Successfully Deleted User!")
}