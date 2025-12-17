// file: config_api.go
package api

import (
	"net/http"
	"strconv"

	"backend/models"

	"github.com/gin-gonic/gin"
)

// GetConfigUnits 获取所有单位配置
func GetConfigUnits(c *gin.Context) {
	units, err := configRepo.GetAllUnits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取单位配置失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, units)
}

// AddConfigUnit 添加单位配置
func AddConfigUnit(c *gin.Context) {
	var unit models.ConfigUnit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if unit.Name == "" || unit.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "单位名称和代码不能为空"})
		return
	}

	if err := configRepo.AddUnit(&unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加单位失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "单位添加成功",
		"unit":    unit,
	})
}

// GetConfigSubstations 获取所有变电站配置
func GetConfigSubstations(c *gin.Context) {
	substations, err := configRepo.GetAllSubstations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取变电站配置失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, substations)
}

// GetSubstationsGrouped 获取按电压等级分组的变电站
func GetSubstationsGrouped(c *gin.Context) {
	grouped, err := configRepo.GetSubstationsByVoltageLevel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取变电站分组失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, grouped)
}

// AddConfigSubstation 添加变电站配置
func AddConfigSubstation(c *gin.Context) {
	var substation models.ConfigSubstation
	if err := c.ShouldBindJSON(&substation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if substation.Name == "" || substation.VoltageLevel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "变电站名称和电压等级不能为空"})
		return
	}

	if substation.Status == "" {
		substation.Status = "active"
	}

	if err := configRepo.AddSubstation(&substation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加变电站失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "变电站添加成功",
		"substation": substation,
	})
}

// GetSafetyTemplates 获取安全措施模板
func GetSafetyTemplates(c *gin.Context) {
	templates, err := configRepo.GetAllSafetyTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取安全措施模板失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// GetRequiredSafetyMeasures 获取必填安全措施
func GetRequiredSafetyMeasures(c *gin.Context) {
	measures, err := configRepo.GetRequiredSafetyMeasures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取必填安全措施失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"required_measures": measures})
}

// GetDeviceTypes 获取设备类型配置
func GetDeviceTypes(c *gin.Context) {
	deviceTypes, err := configRepo.GetAllDeviceTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设备类型配置失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, deviceTypes)
}

// DeleteConfigUnit 删除单位配置
func DeleteConfigUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的单位ID"})
		return
	}

	if err := configRepo.DeleteUnit(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除单位失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "单位删除成功"})
}

// DeleteConfigSubstation 删除变电站配置
func DeleteConfigSubstation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的变电站ID"})
		return
	}

	if err := configRepo.DeleteSubstation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除变电站失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "变电站删除成功"})
}

// UpdateConfigUnit 更新单位配置
func UpdateConfigUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的单位ID"})
		return
	}

	var unit models.ConfigUnit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	unit.ID = id
	if err := configRepo.UpdateUnit(&unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新单位失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "单位更新成功",
		"unit":    unit,
	})
}

// UpdateConfigSubstation 更新变电站配置
func UpdateConfigSubstation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的变电站ID"})
		return
	}

	var substation models.ConfigSubstation
	if err := c.ShouldBindJSON(&substation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	substation.ID = id
	if err := configRepo.UpdateSubstation(&substation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新变电站失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "变电站更新成功",
		"substation": substation,
	})
}
