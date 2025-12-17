// file: admin.go
package api

import (
	"net/http"
	"strconv"
	"time"

	"backend/models"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Position string `json:"position"`
	Unit     string `json:"unit"`
}

type TemplateRequest struct {
	Name           string                `json:"name" binding:"required"`
	Unit           string                `json:"unit"`
	Substations    []string              `json:"substations"`
	WorkContent    string                `json:"work_content" binding:"required"`
	RelatedDevices []models.ReportDevice `json:"related_devices"`
	Tools          []models.ReportTool   `json:"tools"`
	Operation      string                `json:"operation"`
	SafetyMeasures string                `json:"safety_measures"`
}

func GetAllUsers(c *gin.Context) {
	users, err := userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		RealName: req.RealName,
		Phone:    req.Phone,
		Position: req.Position,
		Role:     req.Role,
		Unit:     req.Unit,
	}

	if err := userRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户添加成功"})
}

func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	if err := userRepo.DeleteUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	username := c.Param("username")
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Username: username,
		RealName: req.RealName,
		Phone:    req.Phone,
		Position: req.Position,
		Role:     req.Role,
		Unit:     req.Unit,
	}

	if err := userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户更新成功"})
}

func GetAllTemplates(c *gin.Context) {
	templates, err := templateRepo.GetAllTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func GetTemplateByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	template, err := templateRepo.GetTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func CreateTemplate(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")

	template := &models.Template{
		Name:           req.Name,
		Unit:           req.Unit,
		Substations:    req.Substations,
		WorkContent:    req.WorkContent,
		RelatedDevices: req.RelatedDevices,
		Tools:          req.Tools,
		Operation:      req.Operation,
		SafetyMeasures: req.SafetyMeasures,
		CreatedBy:      userID.(int),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := templateRepo.CreateTemplate(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建模板失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "模板创建成功",
		"template_id": template.ID,
	})
}

func UpdateTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 检查模板是否存在 - 简化处理，不使用变量
	if _, err := templateRepo.GetTemplateByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	template := &models.Template{
		ID:             id,
		Name:           req.Name,
		Unit:           req.Unit,
		Substations:    req.Substations,
		WorkContent:    req.WorkContent,
		RelatedDevices: req.RelatedDevices,
		Tools:          req.Tools,
		Operation:      req.Operation,
		SafetyMeasures: req.SafetyMeasures,
		UpdatedAt:      time.Now(),
	}

	if err := templateRepo.UpdateTemplate(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新模板失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "模板更新成功"})
}

func DeleteTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	if err := templateRepo.DeleteTemplate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "模板删除成功"})
}

// GetSystemStatus 获取系统状态
func GetSystemStatus(c *gin.Context) {
	// 这里可以添加获取系统状态的逻辑
	c.JSON(http.StatusOK, gin.H{
		"status":    "running",
		"timestamp": time.Now().Unix(),
	})
}

// BackupDatabase 备份数据库
func BackupDatabase(c *gin.Context) {
	// 这里可以添加数据库备份逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "数据库备份功能待实现",
	})
}
