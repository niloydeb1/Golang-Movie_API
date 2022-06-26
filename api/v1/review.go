package v1

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/niloydeb1/Golang-Movie_API/api/common"
	"github.com/niloydeb1/Golang-Movie_API/enums"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"strings"
	"time"
)

func ReviewRouter(g *echo.Group) {
	g.GET("/:id", reviewApi{}.GetByID)
	g.GET("", reviewApi{}.Search)
	g.POST("", reviewApi{}.Post)
	g.DELETE("", reviewApi{}.Delete)
}

type reviewApi struct {
}

// GetByID... GetByID Api
// @Summary Review get by id api
// @Description Api for getiing review by id
// @Tags Review
// @Produce json
// @Param id path string true "review id"
// @Success 200 {object} common.ResponseDTO{data=v1.Review{}}
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/reviews/{id} [GET]
func (r reviewApi) GetByID(context echo.Context) error {
	id := context.Param("id")
	if id == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Review id is not provided", "Operation failed")
	}
	data := v1.Review{}.GetByID(id)
	if data.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Review is not found!", "Please provide a valid review id!")
	}
	return common.GenerateSuccessResponse(context, data, nil, "Operation Successful")
}

// Search... Search Api
// @Summary Search api
// @Description Api for searching reviews
// @Tags Review
// @Produce json
// @Param title query string false "movie title keyword"
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {object} common.ResponseDTO{data=[]v1.Review{}}
// @Forbidden 403 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Router /api/v1/reviews [GET]
func (r reviewApi) Search(context echo.Context) error {
	pagination := getPagination(context)
	title := strings.ToLower(context.QueryParam("title"))
	var query bson.M
	var data []v1.Review
	var total int64
	if title != "" {
		reg := ".*" + title + ".*"
		query = bson.M{
			"movie.Title": bson.M{"$regex": primitive.Regex{
				Pattern: reg,
				Options: "i",
			}},
		}
	}
	data, total = v1.Review{}.Search(query, pagination)
	metadata := common.GetPaginationMetadata(pagination.Page, pagination.Limit, total, int64(len(data)))
	uri := strings.Split(context.Request().RequestURI, "?")[0]
	if pagination.Page > 0 {
		metadata.Links = append(metadata.Links, map[string]string{"prev": uri + "?title=" + context.QueryParam("title") + "&page=" + strconv.FormatInt(pagination.Page-1, 10) + "&limit=" + strconv.FormatInt(pagination.Limit, 10)})
	}
	metadata.Links = append(metadata.Links, map[string]string{"self": uri + "?title=" + context.QueryParam("title") + "&page=" + strconv.FormatInt(pagination.Page, 10) + "&limit=" + strconv.FormatInt(pagination.Limit, 10)})
	if (pagination.Page+1)*pagination.Limit < metadata.TotalCount {
		metadata.Links = append(metadata.Links, map[string]string{"next": uri + "?title=" + context.QueryParam("title") + "&page=" + strconv.FormatInt(pagination.Page+1, 10) + "&limit=" + strconv.FormatInt(pagination.Limit, 10)})
	}
	return common.GenerateSuccessResponse(context, data,
		&metadata, "Successful")
}

// Post... Post Api
// @Summary Post review api
// @Description Api for posting review
// @Tags Review
// @Produce json
// @Param Authorization header string true "Insert your access token while posting review" default(Bearer <Add access token here>)
// @Param data body v1.Review true "dto for posting review"
// @Success 200 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/reviews [POST]
func (r reviewApi) Post(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.ID == "" {
		return common.GenerateUnauthorizedResponse(context, "[ERROR]: User not found!", "Please provide valid user information.")
	}
	formData := v1.Review{}
	if err := context.Bind(&formData); err != nil {
		log.Println("Input Error:", err.Error())
		return common.GenerateErrorResponse(context, nil, "Failed to Bind Input!")
	}
	err = formData.Validate()
	if err != nil {
		return common.GenerateErrorResponse(context, "[ERROR]: Invalid data provided.", err.Error())
	}
	formData.ID = uuid.New().String()
	formData.ReviewerEmail = userFromToken.Email
	formData.ReviewerId = userFromToken.ID
	formData.CreatedAt = time.Now().UTC()
	err = v1.Review{}.Store(formData)
	if err != nil {
		return common.GenerateErrorResponse(context, err, err.Error())
	}
	return common.GenerateSuccessResponse(context, "[SUCCESS]: Review posted successfully", nil, "Operation Successful")
}

// Delete... Delete Api
// @Summary Delete review api
// @Description Api for deleting review
// @Tags Review
// @Produce json
// @Param Authorization header string true "Insert your access token while deleting review" default(Bearer <Add access token here>)
// @Param id path string true "review id"
// @Success 200 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/reviews [DELETE]
func (r reviewApi) Delete(context echo.Context) error {
	userFromToken, err := GetUserTokenDtoFromBearerToken(context, v1.Jwt{})
	if err != nil {
		return common.GenerateErrorResponse(context, err.Error(), "Operation Failed!")
	}
	if userFromToken.ID == "" {
		return common.GenerateUnauthorizedResponse(context, "[ERROR]: User not found!", "Please provide valid user information.")
	}
	id := context.Param("id")
	if id == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Review id is not provided", "Operation failed")
	}
	data := v1.Review{}.GetByID(id)
	if data.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Review is not found!", "Please provide a valid review id!")
	}
	if data.ReviewerId != userFromToken.ID && userFromToken.Role != enums.ADMIN && userFromToken.Role != enums.SUPERADMIN {
		return common.GenerateUnauthorizedResponse(context, "[ERROR]: Insufficient permission", "Operation Failed!")
	}
	err = v1.Review{}.Delete(id)
	if err != nil {
		return common.GenerateErrorResponse(context, err, err.Error())
	}
	return common.GenerateSuccessResponse(context, "[SUCCESS]: Review deleted successfully", nil, "Operation Successful")
}
