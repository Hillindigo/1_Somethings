package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 JWT Claims，包含用户ID和角色信息
type CustomClaims struct {
	UserID uint `json:"user_id"`
	Role   int8 `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 根据用户ID和角色生成 JWT Token
// secret: 签名密钥
// expireHours: Token 有效时长（小时）
func GenerateToken(userID uint, role int8, secret string, expireHours int) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并验证 JWT Token，返回自定义 Claims
func ParseToken(tokenString string, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
