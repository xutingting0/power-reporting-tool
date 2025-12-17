// file: auth.go (JWT版本)
package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("登录请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	log.Printf("收到登录请求: username=%s", req.Username)

	user, err := userRepo.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("用户不存在或查询失败: username=%s, error=%v", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	log.Printf("找到用户: %s", user.Username)

	// 密码验证
	if user.Password != req.Password {
		log.Printf("密码不匹配")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	log.Printf("用户 %s 登录成功", user.Username)

	// 生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"role":      user.Role,
		"unit":      user.Unit,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("生成token失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"username":  user.Username,
			"real_name": user.RealName,
			"role":      user.Role,
			"unit":      user.Unit,
			"phone":     user.Phone,
			"position":  user.Position,
			"user_id":   user.ID,
		},
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
