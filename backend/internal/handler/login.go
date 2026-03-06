package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mobileordering/internal/model"
	"mobileordering/internal/model/response"
	"mobileordering/internal/service"
)

type LoginHandler struct {
	wechatSvc *service.WeChatService
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		wechatSvc: service.NewWeChatService(),
	}
}

// Gin 版本登录接口
func (h *LoginHandler) Login(c *gin.Context) {

	// 1️⃣ 解析请求 JSON
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid JSON format",
		})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Code is required",
		})
		return
	}

	// 2️⃣ 调用微信服务
	wxResp, err := h.wechatSvc.GetSessionInfo(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "Login failed",
		})
		return
	}

	// 3️⃣ 生成业务 Token
	token, err := h.wechatSvc.GenerateToken(wxResp.OpenID, wxResp.SessionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "Token generation failed",
		})
		return
	}

	// 只返回纯数据部分（不要包含 message/code）
	loginData := model.LoginResponse{
		Token:  token,
		OpenID: wxResp.OpenID,
		// 注意：不再需要 Message 字段！
	}

	// 使用统一 Success 函数包装
	response.Success(c, loginData)
}