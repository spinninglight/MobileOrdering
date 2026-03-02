package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"mobileordering/internal/config"
	"mobileordering/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// WeChatService 处理微信相关逻辑
type WeChatService struct{}

func NewWeChatService() *WeChatService {
	return &WeChatService{}
}

// GetSessionInfo 调用微信 jscode2session 接口
func (s *WeChatService) GetSessionInfo(code string) (*model.WeChatSessionResp, error) {
	params := url.Values{}
	params.Add("appid", config.Cfg.WxAppID)
	params.Add("secret", config.Cfg.WxAppSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	reqURL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?%s", params.Encode())

	resp, err := httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	var wxResp model.WeChatSessionResp
	if err := json.Unmarshal(body, &wxResp); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	if wxResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat api error: [%d] %s", wxResp.ErrCode, wxResp.ErrMsg)
	}

	return &wxResp, nil
}

// GenerateToken 生成业务 JWT Token
// 在实际生产中，你应该将 session_key 存入 Redis，key 为 openid 或 token，这里仅作演示
func (s *WeChatService) GenerateToken(openID, sessionKey string) (string, error) {
	// 定义 Claims
	claims := jwt.MapClaims{
		"openid": openID,
		// 注意：session_key 通常不直接放入 JWT (因为 JWT 会发给前端)，
		// 而是存在后端 Redis 中，JWT 里只存 openid，后端通过 openid 去 Redis 查 session_key。
		// 但为了演示完整性，这里我们只存 openid 到 Token。
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	// TODO: 此处应将 sessionKey 存入 Redis
	// key := "session:" + openID
	// redis.Set(key, sessionKey, 24*time.Hour)

	// 模拟打印日志，实际生产请移除
	fmt.Printf("[DEBUG] Generated token for OpenID: %s (SessionKey stored internally)\n", openID)

	return tokenString, nil
}
