package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claims JWT载荷
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var (
	// 这里应该从配置文件中读取，这只是示例
	jwtSecret = []byte("your-secret-key-change-in-production")
)

// GenerateJWT 生成JWT token
func GenerateJWT(userID primitive.ObjectID, username, role string, expireDuration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID.Hex(),
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT 解析JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token签名方法错误")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token无效")
}

// SetJWTSecret 设置JWT密钥
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}
