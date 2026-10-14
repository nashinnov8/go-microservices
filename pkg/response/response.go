package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nashinnov8/go-ecommerce/pkg/apperror"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type PaginationMeta struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int32 `json:"total_pages"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Paginated(c *gin.Context, data interface{}, page, limit int32, total int64) {
	var totalPages int32
	if limit > 0 {
		totalPages = int32((total + int64(limit) - 1) / int64(limit))
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})
}

func Error(c *gin.Context, err error) {
	appErr := apperror.FromError(err)
	c.AbortWithStatusJSON(appErr.Status, APIResponse{
		Success: false,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

func ValidationError(c *gin.Context, details interface{}) {
	c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error: &ErrorBody{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request payload or parameters",
			Details: details,
		},
	})
}
