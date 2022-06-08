package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// MetaData Http response metadata
type MetaData struct {
	Page       int64               `json:"page"`
	PerPage    int64               `json:"per_page"`
	PageCount  int64               `json:"page_count"`
	TotalCount int64               `json:"total_count"`
	Links      []map[string]string `json:"links"`
}

// ResponseDTO Http response dto
type ResponseDTO struct {
	Data    interface{} `json:"data" msgpack:"data" xml:"data"`
	Status  string      `json:"status" msgpack:"status" xml:"status"`
	Message string      `json:"message" msgpack:"message" xml:"message"`
}

// ResponseDTOWithPagination Http response dto with pagination
type ResponseDTOWithPagination struct {
	Metadata *MetaData   `json:"_metadata"`
	Data     interface{} `json:"data" msgpack:"data" xml:"data"`
	Status   string      `json:"status" msgpack:"status" xml:"status"`
	Message  string      `json:"message" msgpack:"message" xml:"message"`
}

// GenerateSuccessResponse Http success response
func GenerateSuccessResponse(c echo.Context, data interface{}, metadata *MetaData, message string) error {
	if metadata != nil {
		return c.JSON(http.StatusOK, ResponseDTOWithPagination{
			Status:   "success",
			Message:  message,
			Data:     data,
			Metadata: metadata,
		})
	}
	return c.JSON(http.StatusOK, ResponseDTO{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// GenerateErrorResponse Http error response
func GenerateErrorResponse(c echo.Context, data interface{}, message string) error {
	return c.JSON(http.StatusBadRequest, ResponseDTO{
		Status:  "error",
		Message: message,
		Data:    data,
	})
}

// GenerateForbiddenResponse Http forbidden response
func GenerateForbiddenResponse(c echo.Context, data interface{}, message string) error {
	return c.JSON(http.StatusForbidden, ResponseDTO{
		Status:  "forbidden",
		Message: message,
		Data:    data,
	})
}

// GenerateUnauthorizedResponse Http unauthorized response
func GenerateUnauthorizedResponse(c echo.Context, data interface{}, message string) error {
	return c.JSON(http.StatusUnauthorized, ResponseDTO{
		Status:  "unauthorized",
		Message: message,
		Data:    data,
	})
}

// GetPaginationMetadata return pagination metadata
func GetPaginationMetadata(page, limit, totalRecords, totalPaginatedRecords int64) MetaData {
	metaData := MetaData{
		Page:       page,
		PerPage:    limit,
		TotalCount: totalRecords,
		PageCount:  totalPaginatedRecords,
	}
	return metaData
}
