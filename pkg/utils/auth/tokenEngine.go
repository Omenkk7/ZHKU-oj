// token_engine.go
package auth

import "time"

// TokenEngine 调度器，统一管理 token 实现
type TokenEngine struct {
	impl Token
}

// NewTokenEngine 根据类型和配置创建引擎
func NewTokenEngine(t TokenType, secret string, expiration time.Duration) *TokenEngine {
	return &TokenEngine{
		impl: TokenFactory(t, secret, expiration),
	}
}

// Generate 调用具体实现
func (e *TokenEngine) Generate(userID string, claims map[string]interface{}) (string, error) {
	return e.impl.Generate(userID, claims)
}

// Parse 调用具体实现
func (e *TokenEngine) Parse(tokenString string) (map[string]interface{}, error) {
	return e.impl.Parse(tokenString)
}

// Refresh 调用具体实现
func (e *TokenEngine) Refresh(tokenString string) (string, error) {
	return e.impl.Refresh(tokenString)
}
