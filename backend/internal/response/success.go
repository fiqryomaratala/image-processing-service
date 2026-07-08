package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, message string, data any, meta any) {
	c.JSON(http.StatusOK, SuccessBody{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Created(c *gin.Context, message string, data any, meta any) {
	c.JSON(http.StatusCreated, SuccessBody{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
