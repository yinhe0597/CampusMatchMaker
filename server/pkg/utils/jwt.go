package utils

import (
	"time"

	"campus_collab/internal/infra/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 自定义声明
type Claims struct {
	UserID    uint   `json:"user_id"`
	StudentID string `json:"student_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, studentID string, jwtCfg config.JWTConfig) (string, int64, error) {
	expireAt := time.Now().Add(time.Duration(jwtCfg.ExpireHours) * time.Hour)

	claims := &Claims{
		UserID:    userID,
		StudentID: studentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "campus_collab",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expireAt.Unix(), nil
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
