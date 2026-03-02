package model

// WeChatSessionResp 微信官方接口返回结构
type WeChatSessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// LoginRequest 前端发送的请求体
type LoginRequest struct {
	Code string `json:"code"`
}

// LoginResponse 后端返回给前端的响应
type LoginResponse struct {
	Token   string `json:"token"`   // 业务 Token，前端后续请求携带此 Token
	OpenID  string `json:"openid"`  // 可选：返回 OpenID 给前端做展示用
	Message string `json:"message"`
}

// ErrorResponse 通用错误响应
type ErrorResponse struct {
	Error string `json:"error"`
}