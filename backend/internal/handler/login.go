package handler

import (
	"encoding/json"
	"net/http"

	"mobileordering/internal/model"
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

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 设置 CORS (开发环境可能需要，生产环境建议配置 Nginx)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, model.ErrorResponse{Error: "Method not allowed"})
		return
	}

	// 1. 解析请求
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	if req.Code == "" {
		sendJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "Code is required"})
		return
	}

	// 2. 调用微信服务
	wxResp, err := h.wechatSvc.GetSessionInfo(req.Code)
	if err != nil {
		// 记录详细错误日志，返回通用错误给前端
		// log.Printf("WeChat login failed: %v", err)
		sendJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "Login failed: " + err.Error()})
		return
	}

	// 3. 生成业务 Token
	token, err := h.wechatSvc.GenerateToken(wxResp.OpenID, wxResp.SessionKey)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "Token generation failed"})
		return
	}

	// 4. 返回成功响应
	// 注意：绝对不要返回 session_key 给前端！
	response := model.LoginResponse{
		Token:   token,
		OpenID:  wxResp.OpenID, // 可选返回
		Message: "Login successful",
	}

	sendJSON(w, http.StatusOK, response)
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
