// token_factory.go
package auth

import "time"

// TokenType 枚举不同实现
type TokenType string

const (
	TokenJWT   TokenType = "jwt"
	TokenRedis TokenType = "redis"
)

// TokenFactory 创建 Token 实现
func TokenFactory(t TokenType, secret string, expiration time.Duration) Token {
	switch t {
	case TokenJWT:
		return NewJWTToken(secret, expiration)
	case TokenRedis:
		// TODO: Redis 实现
		panic("Redis token not implemented yet")
	default:
		panic("unsupported token type")
	}
}
