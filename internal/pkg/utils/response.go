package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response - struct untuk response API yang konsisten
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo - struct untuk error detail
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse - kirim response sukses dengan data
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: "success",
		Data:   data,
	})
}

// CreatedResponse - kirim response created (201)
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Status: "success",
		Data:   data,
	})
}

// ErrorResponse - kirim response error
func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Status: "error",
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// ValidationErrorResponse - kirim response error validasi
func ValidationErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", message)
}
