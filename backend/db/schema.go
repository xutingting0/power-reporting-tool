// file: db/schema.go
package db

import (
	"log"
)

func CreateTables() {
	// 用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		real_name TEXT NOT NULL,
		phone TEXT,
		position TEXT,
		role TEXT NOT NULL DEFAULT 'user',
		unit TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 报备单表
	createReportsTable := `
	CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		report_no TEXT NOT NULL UNIQUE,
		reporter_id INTEGER NOT NULL,
		reporter_name TEXT NOT NULL,
		unit TEXT NOT NULL,
		substations TEXT NOT NULL, -- JSON数组存储
		work_time_start DATETIME NOT NULL,
		work_time_end DATETIME NOT NULL,
		responsible_person TEXT NOT NULL,
		work_content TEXT NOT NULL,
		related_devices TEXT NOT NULL, -- JSON数组存储
		tools TEXT NOT NULL, -- JSON数组存储
		operation TEXT NOT NULL,
		safety_measures TEXT NOT NULL,
		remarks TEXT,
		status TEXT NOT NULL DEFAULT 'draft', -- draft, submitted, approved, rejected
		review_comments TEXT,
		reviewer_id INTEGER,
		reviewer_name TEXT,
		reviewed_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (reporter_id) REFERENCES users (id)
	);`

	// 操作日志表
	createLogsTable := `
	CREATE TABLE IF NOT EXISTS report_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		report_id INTEGER NOT NULL,
		action TEXT NOT NULL, -- 操作类型：create, submit, approve, reject, modify
		operator_id INTEGER NOT NULL,
		operator_name TEXT NOT NULL,
		comment TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (report_id) REFERENCES reports (id)
	);`

	// 模板表
	createTemplatesTable := `
	CREATE TABLE IF NOT EXISTS templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		unit TEXT,
		substations TEXT, -- JSON数组存储
		work_content TEXT NOT NULL,
		related_devices TEXT, -- JSON数组存储
		tools TEXT, -- JSON数组存储
		operation TEXT,
		safety_measures TEXT,
		created_by INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 单位配置表
	createUnitsTable := `
	CREATE TABLE IF NOT EXISTS config_units (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		code TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 变电站配置表
	createSubstationsTable := `
	CREATE TABLE IF NOT EXISTS config_substations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		voltage_level TEXT NOT NULL, -- 电压等级：110kV, 220kV等
		region TEXT, -- 所属区域
		description TEXT,
		status TEXT DEFAULT 'active', -- active, inactive
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 设备类型配置表
	createDeviceTypesTable := `
	CREATE TABLE IF NOT EXISTS config_device_types (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		category TEXT NOT NULL, -- 分类：测控设备、保护设备、计算机等
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 安全措施模板表
	createSafetyTemplatesTable := `
	CREATE TABLE IF NOT EXISTS config_safety_templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		is_required BOOLEAN DEFAULT 0, -- 是否强制要求
		template_type TEXT NOT NULL, -- standard, virus, disconnect 等
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []string{
		createUsersTable,
		createReportsTable,
		createLogsTable,
		createTemplatesTable,
		createUnitsTable,
		createSubstationsTable,
		createDeviceTypesTable,
		createSafetyTemplatesTable,
	}

	for _, tableSQL := range tables {
		_, err := DB.Exec(tableSQL)
		if err != nil {
			log.Fatal("创建表失败:", err, "SQL:", tableSQL)
		}
	}
}
