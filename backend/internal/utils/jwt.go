package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// 1. 声明为包级私有变量，不给初始值
var secretKey []byte

// 2. 提供一个 Init 函数，供 main.go 读取配置后调用
func Init(secret string) {
	if secret == "" {
		panic("JWT secret cannot be empty!") // 启动时如果密钥为空直接阻断，防止带病上线
	}
	secretKey = []byte(secret)
}

// CustomClaims 自定义 JWT 载荷
type CustomClaims struct {
	MerchantID int64 `json:"merchant_id"`
	ShopID     int64 `json:"shop_id"`
	jwtv5.RegisteredClaims
}

// GenerateMerchantToken 生成商家 Token (代码基本不变)
func GenerateMerchantToken(merchantID, shopID int64) (string, error) {
	claims := CustomClaims{
		MerchantID: merchantID,
		ShopID:     shopID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
			Issuer:    "mobile-ordering-system",
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	// 这里直接使用 secretKey
	return token.SignedString(secretKey)
}

// ParseMerchantToken 解析和校验商家 Token (代码基本不变)
func ParseMerchantToken(tokenString string) (*CustomClaims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// 这里直接使用 secretKey
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}