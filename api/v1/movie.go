package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/niloydeb1/Golang-Movie_API/api/common"
	v1 "github.com/niloydeb1/Golang-Movie_API/src/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

func MovieRouter(g *echo.Group) {
	g.GET("/:id", movieApi{}.GetByID)
	g.GET("", movieApi{}.Search)
}

type movieApi struct {
}

// GetByID... GetByID Api
// @Summary Movie get by id api
// @Description Api for getiing movie by id
// @Tags Movie
// @Produce json
// @Param id path string true "movie id"
// @Success 200 {object} common.ResponseDTO{data=v1.Movie{}}
// @Failure 400 {object} common.ResponseDTO
// @Forbidden 403 {object} common.ResponseDTO
// @Router /api/v1/movies/{id} [GET]
func (m movieApi) GetByID(context echo.Context) error {
	id := context.Param("id")
	if id == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Movie id is not provided", "Operation failed")
	}
	data := v1.Movie{}.GetByID(id)
	if data.ID == "" {
		return common.GenerateErrorResponse(context, "[ERROR]: Movie is not found!", "Please provide a valid movie id!")
	}
	return common.GenerateSuccessResponse(context, data, nil, "Success!")
}

// Search... Search Api
// @Summary Search api
// @Description Api for searching movies
// @Tags Movie
// @Produce json
// @Param title query string true "title keyword"
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {object} common.ResponseDTO{data=[]v1.Movie{}}
// @Forbidden 403 {object} common.ResponseDTO
// @Failure 400 {object} common.ResponseDTO
// @Router /api/v1/movies [GET]
func (m movieApi) Search(context echo.Context) error {
	pagination := getPagination(context)
	title := strings.ToLower(context.QueryParam("title"))
	query := bson.M{
		"Title": title,
	}
	data, total := v1.Movie{}.Search(query, v1.Pagination{})
	if len(data) == 0 {
		var movie v1.Movie
		_, res, err := v1.HttpClientService{}.Get("https://www.omdbapi.com/?apikey=1154146a&t="+title, nil)
		if err != nil {
			return common.GenerateErrorResponse(context, "[ERROR]: Failed to connect to Omdb server", "Operation failed")
		}
		err = json.Unmarshal(res, &movie)
		if err != nil {
			return common.GenerateErrorResponse(context, err, err.Error())
		}
		if movie.Title == "" {
			return common.GenerateErrorResponse(context, "[ERROR]: Movie does not exist", "Operation failed")
		}
		movie.Title = strings.ToLower(movie.Title)
		checkMovie := v1.Movie{}.GetByTitle(movie.Title)
		if checkMovie.Title == "" {
			movie.ID = uuid.New().String()
			err = v1.Movie{}.Store(movie)
			if err != nil {
				return common.GenerateErrorResponse(context, err, err.Error())
			}
		}
	}
	reg := ".*" + title + ".*"
	query = bson.M{
		"Title": bson.M{"$regex": primitive.Regex{
			Pattern: reg,
			Options: "i",
		}},
	}
	data, total = v1.Movie{}.Search(query, pagination)
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

func getPagination(context echo.Context) v1.Pagination {
	option := v1.Pagination{}
	page := context.QueryParam("page")
	limit := context.QueryParam("limit")
	if page == "" {
		option.Page = 0
		option.Limit = 10
	} else {
		option.Page, _ = strconv.ParseInt(page, 10, 64)
		option.Limit, _ = strconv.ParseInt(limit, 10, 64)
	}
	return option
}
