package v1

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/niloydeb1/Golang-Movie_API/api/common"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"log"
	"time"
)

func CommentRouter(g *echo.Group) {
	g.GET("/:id", commentApi{}.GetByID)
	g.POST("", commentApi{}.Post)
	g.DELETE("/:id", commentApi{}.Delete)
}

type commentApi struct {
}

// GetByID... GetByID Api
// @Summary Comment get by id api
// @Description Api for getiing comment by id
// @Tags Comment
// @Produce json
// @Param id path string true "comment id"
// @Success 200 {object} common.ResponseDTO{data=v1.Comment{}}
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/comments/{id} [GET]
func (c commentApi) GetByID(context echo.Context) error {
	id := context.Param("id")
	if id == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Comment id is not provided", "Operation failed")
	}
	data := v1.Comment{}.GetByID(id)
	if data.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Comment is not found!", "Please provide a valid comment id!")
	}
	return common.GenerateSuccessResponse(context, data, nil, "Operation Successful")
}

// Post... Post Api
// @Summary Post comment api
// @Description Api for posting comment
// @Tags Comment
// @Produce json
// @Param Authorization header string true "Insert your access token while posting comment" default(Bearer <Add access token here>)
// @Param data body v1.Comment true "dto for posting comment"
// @Success 200 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/comments [POST]
func (c commentApi) Post(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.ID == "" {
		return common.GenerateUnauthorizedResponse(context, "[ERROR]: User not found!", "Please provide valid user information.")
	}
	commentDto := v1.Comment{}
	if err := context.Bind(&commentDto); err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, nil, "Failed to Bind Input!")
	}
	err = commentDto.Validate()
	if err != nil {
		return common.GenerateErrorResponse(context, "[ERROR]: Invalid data provided", err.Error())
	}
	review := v1.Review{}.GetByID(commentDto.ReviewId)
	if review.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Review is not found", "Operation Failed")
	}
	commentDto.MovieId = review.Movie.ID
	commentDto.ID = uuid.New().String()
	commentDto.CommenterId = userFromToken.ID
	commentDto.CommenterEmail = userFromToken.Email
	commentDto.CreatedAt = time.Now().UTC()
	err = v1.Comment{}.Store(commentDto)
	if err != nil {
		return common.GenerateErrorResponse(context, err, err.Error())
	}
	return common.GenerateSuccessResponse(context, "[SUCCESS]: Comment is posted successfully", nil, "Operation Successful")
}

// Delete... Delete Api
// @Summary Delete comment api
// @Description Api for deleting comment
// @Tags Comment
// @Produce json
// @Param Authorization header string true "Insert your access token while deleting comment" default(Bearer <Add access token here>)
// @Param id path string true "comment id"
// @Success 200 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/comments/{id} [DELETE]
func (c commentApi) Delete(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.ID == "" {
		return common.GenerateUnauthorizedResponse(context, "[ERROR]: User not found!", "Please provide valid user information.")
	}
	id := context.Param("id")
	if id == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Comment id is not provided", "Operation failed")
	}
	comment := v1.Comment{}.GetByID(id)
	if comment.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Comment is not found!", "Please provide a valid comment id!")
	}
	if comment.CommenterId != userFromToken.ID && userFromToken.Role != enums.ADMIN && userFromToken.Role != enums.SUPERADMIN {
		return common.GenerateUnauthorizedResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	err = v1.Comment{}.Delete(id)
	if err != nil {
		return common.GenerateErrorResponse(context, err, err.Error())
	}
	return common.GenerateSuccessResponse(context, "[SUCCESS]: Comment is deleted successfully", nil, "Operation Successful")
}
