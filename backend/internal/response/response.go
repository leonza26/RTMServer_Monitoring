package response

import "github.com/gin-gonic/gin"

// successResponse adalah format standar untuk response sukses.
type successResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

// errorBody berisi detail error yang dikirim ke client.
type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// errorResponse adalah format standar untuk response error.
type errorResponse struct {
	Error errorBody `json:"error"`
}

// Success mengirim response sukses dengan format {"data": ...}.
func Success(c *gin.Context, status int, data any) {
	c.JSON(status, successResponse{Data: data})
}

// Error mengirim response error dengan format {"error": {"code": ..., "message": ...}}
// lalu menghentikan handler/middleware berikutnya.
func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, errorResponse{
		Error: errorBody{Code: code, Message: message},
	})
}
