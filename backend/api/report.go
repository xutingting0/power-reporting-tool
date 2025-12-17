// file: report.go
package api

import (
	"backend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-this")

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	user, err := models.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  user.Username,
		"real_name": user.RealName,
		"role":      user.Role,
		"unit":      user.Unit,
		"phone":     user.Phone,
		"position":  user.Position,
		"user_id":   user.ID,
	})
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token声明"})
			c.Abort()
			return
		}

		// 验证用户是否存在
		userID := int(claims["user_id"].(float64))
		user, err := models.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", user.ID)
		c.Set("username", user.Username)
		c.Set("realName", user.RealName)
		c.Set("userRole", user.Role)
		c.Set("userUnit", user.Unit)

		c.Next()
	}
}

// AdminMiddleware 管理员中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 辅助函数
func generateToken(userID int, username, role string) string {
	// 简化版token生成，实际项目应使用JWT
	return username + "-" + role + "-" + string(rune(userID+1000))
}

func parseToken(token string) (int, string, string, error) {
	// 简化版token解析
	parts := strings.Split(token, "-")
	if len(parts) < 3 {
		return 0, "", "", gin.Error{}
	}
	username := parts[0]
	role := parts[1]
	userID := 0
	if len(parts) > 2 {
		userID = int(parts[2][0]) - 1000
	}
	return userID, username, role, nil
}
