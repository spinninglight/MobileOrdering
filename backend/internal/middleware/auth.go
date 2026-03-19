// internal/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "mobileordering/internal/utils"
)

// MerchantAuth 商家认证中间件
func MerchantAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取 Authorization 头
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization 请求头"})
            c.Abort()
            return
        }

        // 2. 校验 Bearer 格式
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 格式错误，需为 Bearer {token}"})
            c.Abort()
            return
        }

        // 3. 解析 Token
        claims, err := jwt.ParseMerchantToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期"})
            c.Abort()
            return
        }

        // 4. 将商家信息注入 Context（这里对应了你最开始 GetPendingOrders 里的取值逻辑）
        c.Set("merchant_id", claims.MerchantID)
        c.Set("shop_id", claims.ShopID)

        // 5. 放行
        c.Next()
    }
}