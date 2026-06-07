package middleware

import (
	"strings"

	"campus_collab/internal/infra/config"
	"campus_collab/pkg/response"
	"campus_collab/pkg/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth(jwtCfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization Header 提取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 校验 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := utils.ParseToken(parts[1], jwtCfg.Secret)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 设置用户信息到 Context
		c.Set("user_id", claims.UserID)
		c.Set("student_id", claims.StudentID)
		c.Next()
	}
}
