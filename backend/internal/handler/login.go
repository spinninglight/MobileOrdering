package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mobileordering/internal/model"
	"mobileordering/internal/model/response"
	"mobileordering/internal/service"
	"mobileordering/internal/utils"
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


// 假设你有这个结构体
type MerchantLoginReq struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (h *LoginHandler) MerchantLogin(c *gin.Context) {
    var req MerchantLoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
        return
    }

    // 1. 调用 Service 验证账号密码，获取商户ID和店铺ID
    // TODO: 这里需要你查数据库比对密码
    // merchantID, shopID, err := h.merchantSvc.Verify(req.Username, req.Password)
    
    // 模拟验证成功获取到的数据：
    merchantID := int64(1001) 
    shopID := int64(1001)

    // 2. 使用我们写好的 JWT 组件签发 Token
    tokenString, err := jwt.GenerateMerchantToken(merchantID, shopID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "生成凭证失败"})
        return
    }

    // 3. 返回给前端
    response.Success(c, gin.H{
        "token":   tokenString,
        "shop_id": shopID,
    })
}