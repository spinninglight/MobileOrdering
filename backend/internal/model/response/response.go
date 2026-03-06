// internal/response/response.go
package response

import "github.com/gin-gonic/gin"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"` // omitempty: 如果 data 为 nil，JSON 中不显示该字段
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(getHTTPStatus(code), Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// 辅助函数：根据业务 code 映射 HTTP 状态码（可选）
func getHTTPStatus(code int) int {
	switch code {
	case CodeBadRequest:
		return 400
	case CodeUnauthorized:
		return 401
	case CodeForbidden:
		return 403
	case CodeNotFound:
		return 404
	case CodeServerError:
		return 500
	default:
		return 200
	}
}