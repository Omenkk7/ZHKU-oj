// token.go
package auth

// Token 抽象接口
type Token interface {
	Generate(userID string, claims map[string]interface{}) (string, error)
	Parse(tokenString string) (map[string]interface{}, error)
	Refresh(tokenString string) (string, error)
}
