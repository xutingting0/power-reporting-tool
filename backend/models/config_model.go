// file: models/config_model.go
package models

// ConfigUnit 单位配置
type ConfigUnit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

// ConfigSubstation 变电站配置
type ConfigSubstation struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	VoltageLevel string `json:"voltage_level"`
	Region       string `json:"region,omitempty"`
	Description  string `json:"description,omitempty"`
	Status       string `json:"status"`
}

// ConfigDeviceType 设备类型配置
type ConfigDeviceType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
}

// ConfigSafetyTemplate 安全措施模板
type ConfigSafetyTemplate struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	IsRequired bool   `json:"is_required"`
	Type       string `json:"type"`
	SortOrder  int    `json:"sort_order"`
}
