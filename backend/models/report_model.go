// file: models/report_model.go
package models

import (
	"database/sql"
	"time"
)

// ReportDevice 报备设备
type ReportDevice struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Type string `json:"type,omitempty"`
}

// ReportTool 报备工具
type ReportTool struct {
	Name         string `json:"name"`
	MAC          string `json:"mac,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	Type         string `json:"type,omitempty"`
}

// Report 报备单模型
type Report struct {
	ID                int            `json:"id"`
	ReportNo          string         `json:"report_no"`
	ReporterID        int            `json:"reporter_id"`
	ReporterName      string         `json:"reporter_name"`
	Unit              string         `json:"unit"`
	Substations       []string       `json:"substations"`
	WorkTimeStart     time.Time      `json:"work_time_start"`
	WorkTimeEnd       time.Time      `json:"work_time_end"`
	ResponsiblePerson string         `json:"responsible_person"`
	WorkContent       string         `json:"work_content"`
	RelatedDevices    []ReportDevice `json:"related_devices"`
	Tools             []ReportTool   `json:"tools"`
	Operation         string         `json:"operation"`
	SafetyMeasures    string         `json:"safety_measures"`
	Remarks           string         `json:"remarks"`
	Status            string         `json:"status"`
	ReviewComments    string         `json:"review_comments"`
	ReviewerID        int            `json:"reviewer_id,omitempty"`
	ReviewerName      string         `json:"reviewer_name,omitempty"`
	ReviewedAt        sql.NullTime   `json:"reviewed_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// ReportLog 报备单操作日志
type ReportLog struct {
	ID           int       `json:"id"`
	ReportID     int       `json:"report_id"`
	Action       string    `json:"action"`
	OperatorID   int       `json:"operator_id"`
	OperatorName string    `json:"operator_name"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}

// Template 模板模型
type Template struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Unit           string         `json:"unit"`
	Substations    []string       `json:"substations"`
	WorkContent    string         `json:"work_content"`
	RelatedDevices []ReportDevice `json:"related_devices"`
	Tools          []ReportTool   `json:"tools"`
	Operation      string         `json:"operation"`
	SafetyMeasures string         `json:"safety_measures"`
	CreatedBy      int            `json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
