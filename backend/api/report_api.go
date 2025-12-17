package api

import (
	"backend/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// CreateReport 创建报备单
func CreateReport(c *gin.Context) {
	var report models.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证安全措施是否包含必填项
	if !validateSafetyMeasures(report.SafetyMeasures) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "安全措施必须包含'已确认所有工具完成病毒扫描'"})
		return
	}

	// 从上下文中获取用户信息
	userID, _ := c.Get("userID")
	userName, _ := c.Get("realName")
	userUnit, _ := c.Get("userUnit")

	report.ReporterID = userID.(int)
	report.ReporterName = userName.(string)
	report.Unit = userUnit.(string)

	// 生成报备单号
	report.ReportNo = generateReportNo()

	// 设置默认状态
	if report.Status == "" {
		report.Status = "draft"
	}

	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	if err := reportRepo.CreateReport(&report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建报备单失败"})
		return
	}

	// 记录操作日志
	_ = reportRepo.CreateReportLog(&models.ReportLog{
		ReportID:     report.ID,
		Action:       "create",
		OperatorID:   userID.(int),
		OperatorName: userName.(string),
		Comment:      "创建报备单",
		CreatedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message":   "报备单创建成功",
		"report_id": report.ID,
		"report_no": report.ReportNo,
	})
}

// GetUserReports 获取用户报备单
func GetUserReports(c *gin.Context) {
	userID, _ := c.Get("userID")

	// 获取查询参数
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	reportNo := c.Query("report_no")
	substation := c.Query("substation")

	reports, err := reportRepo.GetReportsByUserWithFilter(userID.(int), status, startDate, endDate, reportNo, substation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取报备单失败"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetPendingReports 获取待审核报备单
func GetPendingReports(c *gin.Context) {
	reports, err := reportRepo.GetPendingReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取待审核报备单失败"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetAllReports 获取所有报备单（管理员）
func GetAllReports(c *gin.Context) {
	// 获取查询参数
	reporterName := c.Query("reporter_name")
	reportNo := c.Query("report_no")
	substation := c.Query("substation")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	reports, err := reportRepo.GetAllReportsWithFilter(reporterName, reportNo, substation, startDate, endDate, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取报备单失败"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetReportByID 获取报备单详情
func GetReportByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报备单ID"})
		return
	}

	report, err := reportRepo.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报备单不存在"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// UpdateReport 更新报备单
func UpdateReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报备单ID"})
		return
	}

	var report models.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证安全措施是否包含必填项
	if !validateSafetyMeasures(report.SafetyMeasures) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "安全措施必须包含'已确认所有工具完成病毒扫描'"})
		return
	}

	report.ID = id
	report.UpdatedAt = time.Now()

	// 获取用户信息
	userID, _ := c.Get("userID")
	userName, _ := c.Get("realName")

	// 检查权限（只能修改自己的草稿或驳回的报备单）
	existingReport, err := reportRepo.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报备单不存在"})
		return
	}

	if existingReport.ReporterID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此报备单"})
		return
	}

	if existingReport.Status != "draft" && existingReport.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能修改草稿或被驳回的报备单"})
		return
	}

	if err := reportRepo.UpdateReport(&report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新报备单失败"})
		return
	}

	// 记录操作日志
	_ = reportRepo.CreateReportLog(&models.ReportLog{
		ReportID:     id,
		Action:       "update",
		OperatorID:   userID.(int),
		OperatorName: userName.(string),
		Comment:      "修改报备单",
		CreatedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "报备单更新成功"})
}

// SubmitReport 提交报备单
func SubmitReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报备单ID"})
		return
	}

	// 获取用户信息
	userID, _ := c.Get("userID")
	userName, _ := c.Get("realName")

	// 检查报备单是否存在且属于当前用户
	report, err := reportRepo.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报备单不存在"})
		return
	}

	if report.ReporterID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权提交此报备单"})
		return
	}

	if report.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能提交草稿状态的报备单"})
		return
	}

	// 更新状态为已提交
	if err := reportRepo.UpdateReportStatus(id, "submitted"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交报备单失败"})
		return
	}

	// 记录操作日志
	_ = reportRepo.CreateReportLog(&models.ReportLog{
		ReportID:     id,
		Action:       "submit",
		OperatorID:   userID.(int),
		OperatorName: userName.(string),
		Comment:      "提交报备单",
		CreatedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "报备单提交成功"})
}

// ApproveReport 审核通过报备单
func ApproveReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报备单ID"})
		return
	}

	var req struct {
		ReviewComments string `json:"review_comments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取审核人信息
	reviewerID, _ := c.Get("userID")
	reviewerName, _ := c.Get("realName")

	// 检查报备单是否存在且状态为已提交
	report, err := reportRepo.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报备单不存在"})
		return
	}

	if report.Status != "submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能审核已提交的报备单"})
		return
	}

	// 更新审核信息 - 使用正确的参数调用
	// reportRepo.ApproveReport 需要 5 个参数: (reportID int, isApproved bool, reviewComments string, reviewerID int, reviewerName string)
	if err := reportRepo.ApproveReport(id, true, req.ReviewComments, reviewerID.(int), reviewerName.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "审核报备单失败"})
		return
	}

	// 记录操作日志
	_ = reportRepo.CreateReportLog(&models.ReportLog{
		ReportID:     id,
		Action:       "approve",
		OperatorID:   reviewerID.(int),
		OperatorName: reviewerName.(string),
		Comment:      req.ReviewComments,
		CreatedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "报备单审核通过"})
}

// RejectReport 驳回报备单
func RejectReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报备单ID"})
		return
	}

	var req struct {
		ReviewComments string `json:"review_comments" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "审核意见不能为空"})
		return
	}

	// 获取审核人信息
	reviewerID, _ := c.Get("userID")
	reviewerName, _ := c.Get("realName")

	// 检查报备单是否存在且状态为已提交
	report, err := reportRepo.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报备单不存在"})
		return
	}

	if report.Status != "submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能驳回已提交的报备单"})
		return
	}

	// 更新驳回信息
	if err := reportRepo.RejectReport(id, req.ReviewComments, reviewerID.(int), reviewerName.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "驳回报备单失败"})
		return
	}

	// 记录操作日志
	_ = reportRepo.CreateReportLog(&models.ReportLog{
		ReportID:     id,
		Action:       "reject",
		OperatorID:   reviewerID.(int),
		OperatorName: reviewerName.(string),
		Comment:      req.ReviewComments,
		CreatedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "报备单已驳回"})
}

// DeleteReport 删除报备单
func DeleteReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报备单ID"})
		return
	}

	// 获取用户信息
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	// 检查报备单是否存在
	report, err := reportRepo.GetReportByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报备单不存在"})
		return
	}

	// 权限检查：管理员可以删除任何报备单，普通用户只能删除自己的草稿
	if userRole.(string) != "admin" && (report.ReporterID != userID.(int) || report.Status != "draft") {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此报备单"})
		return
	}

	if err := reportRepo.DeleteReport(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除报备单失败"})
		return
	}

	// 记录操作日志
	_ = reportRepo.CreateReportLog(&models.ReportLog{
		ReportID:     id,
		Action:       "delete",
		OperatorID:   userID.(int),
		OperatorName: c.GetString("realName"),
		Comment:      "删除报备单",
		CreatedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "报备单删除成功"})
}

// BatchDeleteReports 批量删除报备单
func BatchDeleteReports(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	// 检查权限并批量删除
	for _, id := range req.IDs {
		report, err := reportRepo.GetReportByID(id)
		if err != nil {
			continue
		}

		if userRole.(string) != "admin" && (report.ReporterID != userID.(int) || report.Status != "draft") {
			continue
		}

		_ = reportRepo.DeleteReport(id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量删除成功"})
}

// GetReportLogs 获取报备单操作日志
func GetReportLogs(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报备单ID"})
		return
	}

	logs, err := reportRepo.GetReportLogs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取操作日志失败"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// ExportReports 导出报备单
func ExportReports(c *gin.Context) {
	// 获取查询参数
	reporterName := c.Query("reporter_name")
	reportNo := c.Query("report_no")
	substation := c.Query("substation")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	reports, err := reportRepo.GetAllReportsWithFilter(reporterName, reportNo, substation, startDate, endDate, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取报备单失败"})
		return
	}

	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := "报备单列表"
	f.NewSheet(sheetName)

	// 设置表头
	headers := []string{"报备单号", "报备人", "单位", "变电站", "工作开始时间",
		"工作结束时间", "工作内容", "状态", "创建时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 填充数据
	for i, report := range reports {
		row := i + 2

		// 变电站转为字符串
		substationStr := ""
		if len(report.Substations) > 0 {
			substationStr = report.Substations[0]
			if len(report.Substations) > 1 {
				substationStr += "等"
			}
		}

		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), report.ReportNo)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), report.ReporterName)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), report.Unit)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), substationStr)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), report.WorkTimeStart.Format("2006-01-02 15:04"))
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), report.WorkTimeEnd.Format("2006-01-02 15:04"))
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), report.WorkContent)
		f.SetCellValue(sheetName, "H"+strconv.Itoa(row), getStatusText(report.Status))
		f.SetCellValue(sheetName, "I"+strconv.Itoa(row), report.CreatedAt.Format("2006-01-02 15:04"))
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 20)
	f.SetColWidth(sheetName, "E", "F", 18)
	f.SetColWidth(sheetName, "G", "G", 30)
	f.SetColWidth(sheetName, "H", "H", 10)
	f.SetColWidth(sheetName, "I", "I", 18)

	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=reports_"+time.Now().Format("20060102150405")+".xlsx")

	// 写入响应
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Excel文件失败"})
		return
	}
}

func getStatusText(status string) string {
	statusMap := map[string]string{
		"draft":     "草稿",
		"submitted": "已提交",
		"approved":  "已通过",
		"rejected":  "已驳回",
	}
	if text, ok := statusMap[status]; ok {
		return text
	}
	return status
}

// GetReportStatistics 获取报备单统计
func GetReportStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := reportRepo.GetReportStatistics(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// 辅助函数：生成报备单号
func generateReportNo() string {
	return "REP" + time.Now().Format("20060102150405")
}

// UpdateUserProfile 更新用户个人信息
func UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var req struct {
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Position string `json:"position"`
		Unit     string `json:"unit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户
	user, err := reportRepo.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 更新用户信息
	user.RealName = req.RealName
	user.Phone = req.Phone
	user.Position = req.Position
	user.Unit = req.Unit

	if err := reportRepo.UpdateUserProfile(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}

// GetDashboardStats 获取仪表盘统计
func GetDashboardStats(c *gin.Context) {
	stats, err := reportRepo.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRecentReports 获取最近报备单
func GetRecentReports(c *gin.Context) {
	limit := 10
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	reports, err := reportRepo.GetRecentReports(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取报备单失败"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// UseTemplate 使用模板创建报备单
func UseTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	template, err := reportRepo.GetTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 从上下文中获取用户信息
	userID, _ := c.Get("userID")
	userName, _ := c.Get("realName")
	// userUnit变量未使用，从模板获取单位
	// userUnit, _ := c.Get("userUnit")

	report := &models.Report{
		ReporterID:     userID.(int),
		ReporterName:   userName.(string),
		Unit:           template.Unit, // 使用模板的单位
		Substations:    template.Substations,
		WorkContent:    template.WorkContent,
		RelatedDevices: template.RelatedDevices,
		Tools:          template.Tools,
		Operation:      template.Operation,
		SafetyMeasures: template.SafetyMeasures,
		Status:         "draft",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 如果模板没有单位，使用用户的单位
	if report.Unit == "" {
		if userUnit, exists := c.Get("userUnit"); exists {
			report.Unit = userUnit.(string)
		}
	}

	// 设置默认的工作时间（当前时间到2小时后）
	report.WorkTimeStart = time.Now()
	report.WorkTimeEnd = time.Now().Add(2 * time.Hour)
	report.ResponsiblePerson = report.ReporterName

	// 生成报备单号
	report.ReportNo = generateReportNo()

	if err := reportRepo.CreateReport(report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建报备单失败: " + err.Error()})
		return
	}

	// 记录操作日志
	_ = reportRepo.CreateReportLog(&models.ReportLog{
		ReportID:     report.ID,
		Action:       "create",
		OperatorID:   userID.(int),
		OperatorName: userName.(string),
		Comment:      "使用模板创建报备单",
		CreatedAt:    time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message":   "使用模板创建成功",
		"report_id": report.ID,
		"report_no": report.ReportNo,
	})
}

// 安全措施验证函数
func validateSafetyMeasures(measures string) bool {
	// 必须包含病毒扫描确认
	requiredPhrases := []string{
		"已确认所有工具完成病毒扫描",
		"病毒扫描完成",
		"完成病毒扫描",
	}

	measuresLower := strings.ToLower(measures)
	for _, phrase := range requiredPhrases {
		if strings.Contains(measuresLower, strings.ToLower(phrase)) {
			return true
		}
	}

	return false
}
