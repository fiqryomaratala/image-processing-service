package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context, message string, errs []ErrorItem) {
	writeError(c, http.StatusBadRequest, message, errs)
}

func Unauthorized(c *gin.Context, message string, errs []ErrorItem) {
	writeError(c, http.StatusUnauthorized, message, errs)
}

func Forbidden(c *gin.Context, message string, errs []ErrorItem) {
	writeError(c, http.StatusForbidden, message, errs)
}

func NotFound(c *gin.Context, message string, errs []ErrorItem) {
	writeError(c, http.StatusNotFound, message, errs)
}

func Conflict(c *gin.Context, message string, errs []ErrorItem) {
	writeError(c, http.StatusConflict, message, errs)
}

func UnprocessableEntity(c *gin.Context, message string, errs []ErrorItem) {
	writeError(c, http.StatusUnprocessableEntity, message, errs)
}

func InternalServerError(c *gin.Context, message string, errs []ErrorItem) {
	writeError(c, http.StatusInternalServerError, message, errs)
}

func writeError(c *gin.Context, statusCode int, message string, errs []ErrorItem) {
	c.AbortWithStatusJSON(statusCode, ErrorBody{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}
