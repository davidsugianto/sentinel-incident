package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func DataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Data: data,
		Code: http.StatusOK,
	})
}

func ErrResponse(c *gin.Context, err error, code int) *Response {
	resp := Response{
		Code:  code,
		Error: err.Error(),
	}

	c.AbortWithStatusJSON(code, resp)
	return &resp
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"code": 200, "data": data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(201, gin.H{"code": 201, "data": data})
}

func Fail(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"code": status, "error": err.Error()})
}

func BadRequest(c *gin.Context, err error) {
	Fail(c, 400, err)
}
