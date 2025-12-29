package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("bli-tf-secret-key-change-in-production")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

// ValidateToken 验证 JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

// AuthMiddleware 认证中间件
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		_, err = ValidateToken(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// AuthAPIMiddleware API 认证中间件
func AuthAPIMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, `{"error":"未授权"}`, http.StatusUnauthorized)
			return
		}
		_, err = ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, `{"error":"登录已过期"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
