package apiutil

import "github.com/gin-gonic/gin"

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

func Message(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"message": message})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, errorResponse{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}

func AbortError(c *gin.Context, status int, code, message string) {
	Error(c, status, code, message)
	c.Abort()
}

func NoContent(c *gin.Context) {
	c.Status(204)
}
