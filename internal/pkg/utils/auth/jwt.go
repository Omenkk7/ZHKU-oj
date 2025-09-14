// jwt.go
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTToken 实现 Token 接口
type JWTToken struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTToken(secret string, exp time.Duration) *JWTToken {
	return &JWTToken{
		secret:     []byte(secret),
		expiration: exp,
	}
}

func (j *JWTToken) Generate(userID string, claims map[string]interface{}) (string, error) {
	c := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.expiration).Unix(),
	}
	for k, v := range claims {
		c[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.secret)
}

func (j *JWTToken) Parse(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, err
}

func (j *JWTToken) Refresh(tokenString string) (string, error) {
	claims, err := j.Parse(tokenString)
	if err != nil {
		return "", err
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", err
	}
	return j.Generate(userID, claims)
}
