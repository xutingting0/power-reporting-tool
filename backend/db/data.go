// file: db/data.go
package db

import (
	"log"
)

func InitData() {
	// 初始化用户数据
	initUsers()
	// 初始化配置数据
	initConfigData()
	// 初始化模板数据
	initTemplateData()
}

func initUsers() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Printf("检查用户表失败: %v", err)
		return
	}

	if count == 0 {
		log.Println("初始化用户数据...")

		// 添加管理员账户
		_, err = DB.Exec(`
			INSERT INTO users (username, password, real_name, phone, position, role, unit) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			"admin", "admin123", "系统管理员", "13800000001", "系统管理员", "admin", "遵义供电局")
		if err != nil {
			log.Printf("初始化管理员账户失败: %v", err)
		}

		// 添加普通用户账户
		users := []struct {
			username string
			password string
			realName string
			phone    string
			position string
			unit     string
		}{
			{"user1", "123456", "张三", "13800000002", "运维工程师", "遵义供电局"},
			{"user2", "123456", "李四", "13800000003", "技术员", "贵阳供电局"},
			{"user3", "123456", "王五", "13800000004", "安全员", "六盘水供电局"},
			{"reviewer", "123456", "审核员", "13800000005", "安全审核员", "遵义供电局"},
		}

		for _, user := range users {
			_, err = DB.Exec(`
				INSERT INTO users (username, password, real_name, phone, position, role, unit) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				user.username, user.password, user.realName, user.phone, user.position, "user", user.unit)
			if err != nil {
				log.Printf("初始化用户 %s 失败: %v", user.username, err)
			}
		}

		log.Println("用户数据初始化完成")
	}
}

func initConfigData() {
	log.Println("初始化配置数据...")

	// 1. 初始化单位数据
	var unitCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM config_units").Scan(&unitCount)
	if err != nil {
		log.Printf("检查单位配置表失败: %v", err)
		return
	}

	if unitCount == 0 {
		units := []struct {
			name        string
			code        string
			description string
		}{
			{"遵义供电局", "ZYGD", "遵义地区供电局"},
			{"贵阳供电局", "GYGD", "贵阳地区供电局"},
			{"六盘水供电局", "LPSGD", "六盘水地区供电局"},
			{"安顺供电局", "ASGD", "安顺地区供电局"},
			{"毕节供电局", "BJGD", "毕节地区供电局"},
			{"铜仁供电局", "TRGD", "铜仁地区供电局"},
			{"黔东南供电局", "QDN", "黔东南地区供电局"},
			{"黔西南供电局", "QXN", "黔西南地区供电局"},
		}

		for _, unit := range units {
			_, err = DB.Exec("INSERT INTO config_units (name, code, description) VALUES (?, ?, ?)",
				unit.name, unit.code, unit.description)
			if err != nil {
				log.Printf("初始化单位 %s 失败: %v", unit.name, err)
			}
		}
		log.Println("单位数据初始化完成")
	}

	// 2. 初始化变电站数据
	var substationCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM config_substations").Scan(&substationCount)
	if err != nil {
		log.Printf("检查变电站配置表失败: %v", err)
		return
	}

	if substationCount == 0 {
		substations := []struct {
			name         string
			voltageLevel string
			region       string
			description  string
		}{
			// 110kV变电站
			{"110kV道真变电站", "110kV", "遵义", "道真县110kV变电站"},
			{"110kV正安变电站", "110kV", "遵义", "正安县110kV变电站"},
			{"110kV务川变电站", "110kV", "遵义", "务川县110kV变电站"},
			{"110kV凤冈变电站", "110kV", "遵义", "凤冈县110kV变电站"},
			{"110kV湄潭变电站", "110kV", "遵义", "湄潭县110kV变电站"},
			{"110kV余庆变电站", "110kV", "遵义", "余庆县110kV变电站"},
			{"110kV习水变电站", "110kV", "遵义", "习水县110kV变电站"},

			// 220kV变电站
			{"220kV遵义变电站", "220kV", "遵义", "遵义市220kV枢纽变电站"},
			{"220kV桐梓变电站", "220kV", "遵义", "桐梓县220kV变电站"},
			{"220kV仁怀变电站", "220kV", "遵义", "仁怀市220kV变电站"},
			{"220kV赤水变电站", "220kV", "遵义", "赤水市220kV变电站"},
			{"220kV绥阳变电站", "220kV", "遵义", "绥阳县220kV变电站"},

			// 贵阳地区变电站
			{"220kV贵阳变电站", "220kV", "贵阳", "贵阳市220kV枢纽变电站"},
			{"110kV花溪变电站", "110kV", "贵阳", "花溪区110kV变电站"},
			{"110kV乌当变电站", "110kV", "贵阳", "乌当区110kV变电站"},
			{"110kV白云变电站", "110kV", "贵阳", "白云区110kV变电站"},
		}

		for _, substation := range substations {
			_, err = DB.Exec(`
				INSERT INTO config_substations (name, voltage_level, region, description, status) 
				VALUES (?, ?, ?, ?, 'active')`,
				substation.name, substation.voltageLevel, substation.region, substation.description)
			if err != nil {
				log.Printf("初始化变电站 %s 失败: %v", substation.name, err)
			}
		}
		log.Println("变电站数据初始化完成")
	}

	// 3. 初始化设备类型数据
	var deviceTypeCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM config_device_types").Scan(&deviceTypeCount)
	if err != nil {
		log.Printf("检查设备类型配置表失败: %v", err)
		return
	}

	if deviceTypeCount == 0 {
		deviceTypes := []struct {
			name        string
			category    string
			description string
		}{
			{"测控装置", "测控设备", "用于测量和控制电力系统的装置"},
			{"保护装置", "保护设备", "电力系统继电保护装置"},
			{"后台监控机", "计算机", "监控系统后台计算机"},
			{"远动机", "通讯设备", "远程通信和数据处理设备"},
			{"交换机", "网络设备", "网络数据交换设备"},
			{"路由器", "网络设备", "网络路由设备"},
			{"防火墙", "安全设备", "网络安全防护设备"},
			{"服务器", "服务器", "应用和数据服务器"},
			{"智能终端", "终端设备", "智能操作终端"},
			{"GPS对时装置", "时间设备", "时间同步装置"},
			{"UPS电源", "电源设备", "不间断电源系统"},
			{"录波装置", "记录设备", "故障录波装置"},
		}

		for _, deviceType := range deviceTypes {
			_, err = DB.Exec(`
				INSERT INTO config_device_types (name, category, description) 
				VALUES (?, ?, ?)`,
				deviceType.name, deviceType.category, deviceType.description)
			if err != nil {
				log.Printf("初始化设备类型 %s 失败: %v", deviceType.name, err)
			}
		}
		log.Println("设备类型数据初始化完成")
	}

	// 4. 初始化安全措施模板数据
	var safetyTemplateCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM config_safety_templates").Scan(&safetyTemplateCount)
	if err != nil {
		log.Printf("检查安全措施模板表失败: %v", err)
		return
	}

	if safetyTemplateCount == 0 {
		safetyTemplates := []struct {
			name         string
			content      string
			isRequired   bool
			templateType string
			sortOrder    int
		}{
			{
				name:         "标准安全措施",
				content:      "1. 已确认所有工具完成病毒扫描，无病毒、木马等安全隐患\n2. 已断开非工作设备网络连接\n3. 操作人员已接受安全培训\n4. 已准备应急预案\n5. 已检查安全工器具完好",
				isRequired:   true,
				templateType: "standard",
				sortOrder:    1,
			},
			{
				name:         "病毒扫描确认",
				content:      "已确认所有工具完成病毒扫描，无病毒、木马等安全隐患",
				isRequired:   true,
				templateType: "virus",
				sortOrder:    2,
			},
			{
				name:         "网络连接断开",
				content:      "已断开非工作设备网络连接，仅保留工作必须的网络通道",
				isRequired:   true,
				templateType: "disconnect",
				sortOrder:    3,
			},
			{
				name:         "安全培训完成",
				content:      "操作人员已接受安全培训，了解操作流程和应急预案",
				isRequired:   false,
				templateType: "training",
				sortOrder:    4,
			},
			{
				name:         "工具检查完成",
				content:      "已检查所有安全工器具完好有效，符合安全要求",
				isRequired:   false,
				templateType: "tools",
				sortOrder:    5,
			},
			{
				name:         "应急预案准备",
				content:      "已准备应急预案，应急联系人和联系方式已确认",
				isRequired:   false,
				templateType: "emergency",
				sortOrder:    6,
			},
		}

		for _, template := range safetyTemplates {
			isRequired := 0
			if template.isRequired {
				isRequired = 1
			}

			_, err = DB.Exec(`
				INSERT INTO config_safety_templates (name, content, is_required, template_type, sort_order) 
				VALUES (?, ?, ?, ?, ?)`,
				template.name, template.content, isRequired, template.templateType, template.sortOrder)
			if err != nil {
				log.Printf("初始化安全措施模板 %s 失败: %v", template.name, err)
			}
		}
		log.Println("安全措施模板数据初始化完成")
	}

	log.Println("配置数据初始化完成")
}

func initTemplateData() {
	log.Println("初始化模板数据...")

	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM templates").Scan(&count)
	if err != nil {
		log.Printf("检查模板表失败: %v", err)
		return
	}

	if count == 0 {
		templates := []struct {
			name           string
			unit           string
			substations    string
			workContent    string
			relatedDevices string
			tools          string
			operation      string
			safetyMeasures string
		}{
			{
				name:           "测控装置调试模板",
				unit:           "遵义供电局",
				substations:    `["110kV道真变电站"]`,
				workContent:    "110kV道杠线测控安装及调试",
				relatedDevices: `[{"name":"测控装置","ip":"198.120.0.102","type":"测控设备"},{"name":"后台监控机","ip":"198.120.0.181","type":"计算机"},{"name":"远动机","ip":"123.123.0.1","type":"通讯设备"}]`,
				tools:          `[{"name":"专用运维笔记本电脑","mac":"00-1A-2B-3C-4D-5E","serial_number":"NB-2023-001","type":"计算机"},{"name":"专用安全U盘","serial_number":"USB-2025-ABCDEF","type":"存储设备"},{"name":"调试软件","serial_number":"SOFT-2023-V1.0","type":"软件"}]`,
				operation:      "进行110kV道杠线测控安装、参数配置、功能测试及数据调试",
				safetyMeasures: "1. 已确认所有工具完成病毒扫描，无病毒、木马等安全隐患\n2. 已断开非工作设备网络连接\n3. 操作人员已接受安全培训",
			},
			{
				name:           "设备定期巡检模板",
				unit:           "贵阳供电局",
				substations:    `["220kV贵阳变电站"]`,
				workContent:    "变电站设备定期巡检",
				relatedDevices: `[{"name":"后台监控机","ip":"198.120.0.181","type":"计算机"},{"name":"交换机","ip":"198.120.0.200","type":"网络设备"},{"name":"保护装置","ip":"198.120.0.105","type":"保护设备"}]`,
				tools:          `[{"name":"专用运维笔记本电脑","mac":"00-1A-2B-3C-4D-5F","serial_number":"NB-2023-002","type":"计算机"},{"name":"红外测温仪","serial_number":"IRT-2023-001","type":"仪器"}]`,
				operation:      "1. 检查设备运行状态\n2. 记录设备参数\n3. 检查设备指示灯\n4. 测试设备功能\n5. 清理设备灰尘",
				safetyMeasures: "1. 已确认所有工具完成病毒扫描\n2. 已佩戴安全帽和安全带\n3. 已设置安全警示标志",
			},
			{
				name:           "保护装置定值修改模板",
				unit:           "六盘水供电局",
				substations:    `["110kV盘县变电站"]`,
				workContent:    "保护装置定值修改及验证",
				relatedDevices: `[{"name":"保护装置","ip":"198.120.0.105","type":"保护设备"},{"name":"后台监控机","ip":"198.120.0.182","type":"计算机"}]`,
				tools:          `[{"name":"专用运维笔记本电脑","mac":"00-1A-2B-3C-4D-5G","serial_number":"NB-2023-003","type":"计算机"},{"name":"保护调试仪","serial_number":"PT-2023-001","type":"仪器"}]`,
				operation:      "1. 备份原定值参数\n2. 输入新定值参数\n3. 验证定值正确性\n4. 进行保护功能测试\n5. 记录修改过程",
				safetyMeasures: "1. 已确认所有工具完成病毒扫描\n2. 已断开非必要通信连接\n3. 已准备应急预案\n4. 操作人员双人监护",
			},
			{
				name:           "系统软件升级模板",
				unit:           "安顺供电局",
				substations:    `["220kV安顺变电站"]`,
				workContent:    "监控系统软件升级",
				relatedDevices: `[{"name":"服务器","ip":"198.120.0.190","type":"服务器"},{"name":"后台监控机","ip":"198.120.0.183","type":"计算机"},{"name":"交换机","ip":"198.120.0.201","type":"网络设备"}]`,
				tools:          `[{"name":"专用运维笔记本电脑","mac":"00-1A-2B-3C-4D-5H","serial_number":"NB-2023-004","type":"计算机"},{"name":"专用安全U盘","serial_number":"USB-2025-GHIJK","type":"存储设备"}]`,
				operation:      "1. 备份原系统数据和配置文件\n2. 停止相关服务\n3. 安装升级包\n4. 验证升级结果\n5. 恢复服务运行\n6. 测试系统功能",
				safetyMeasures: "1. 已确认所有工具完成病毒扫描\n2. 已进行系统备份\n3. 已准备回退方案\n4. 操作在业务低谷期进行",
			},
		}

		// 假设管理员用户ID为1
		adminID := 1

		for _, t := range templates {
			_, err = DB.Exec(`
				INSERT INTO templates (name, unit, substations, work_content, related_devices, tools, operation, safety_measures, created_by) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				t.name, t.unit, t.substations, t.workContent, t.relatedDevices, t.tools, t.operation, t.safetyMeasures, adminID)
			if err != nil {
				log.Printf("初始化模板 %s 失败: %v", t.name, err)
			}
		}

		log.Println("模板数据初始化完成，共创建 %d 个模板", len(templates))
	}
}
