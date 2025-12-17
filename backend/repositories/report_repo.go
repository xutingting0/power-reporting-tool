// file: repositories/report_repo.go
package repositories

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) CreateReport(report *models.Report) error {
	substationsJSON, _ := json.Marshal(report.Substations)
	devicesJSON, _ := json.Marshal(report.RelatedDevices)
	toolsJSON, _ := json.Marshal(report.Tools)

	_, err := r.db.Exec(`
		INSERT INTO reports (
			report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			related_devices, tools, operation, safety_measures, remarks, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		report.ReportNo, report.ReporterID, report.ReporterName, report.Unit,
		string(substationsJSON), report.WorkTimeStart, report.WorkTimeEnd,
		report.ResponsiblePerson, report.WorkContent, string(devicesJSON),
		string(toolsJSON), report.Operation, report.SafetyMeasures,
		report.Remarks, report.Status,
	)
	return err
}

func (r *ReportRepository) GetReportByID(id int) (*models.Report, error) {
	var report models.Report
	var substationsJSON, devicesJSON, toolsJSON string
	var reviewedAt sql.NullTime

	query := `
		SELECT id, report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			related_devices, tools, operation, safety_measures, remarks,
			status, review_comments, reviewer_id, reviewer_name, reviewed_at,
			created_at, updated_at
		FROM reports 
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&report.ID, &report.ReportNo, &report.ReporterID, &report.ReporterName,
		&report.Unit, &substationsJSON, &report.WorkTimeStart, &report.WorkTimeEnd,
		&report.ResponsiblePerson, &report.WorkContent, &devicesJSON, &toolsJSON,
		&report.Operation, &report.SafetyMeasures, &report.Remarks, &report.Status,
		&report.ReviewComments, &report.ReviewerID, &report.ReviewerName,
		&reviewedAt, &report.CreatedAt, &report.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 解析JSON字段
	json.Unmarshal([]byte(substationsJSON), &report.Substations)
	json.Unmarshal([]byte(devicesJSON), &report.RelatedDevices)
	json.Unmarshal([]byte(toolsJSON), &report.Tools)

	if reviewedAt.Valid {
		report.ReviewedAt = reviewedAt
	}

	return &report, nil
}

func (r *ReportRepository) UpdateReport(report *models.Report) error {
	substationsJSON, _ := json.Marshal(report.Substations)
	devicesJSON, _ := json.Marshal(report.RelatedDevices)
	toolsJSON, _ := json.Marshal(report.Tools)

	_, err := r.db.Exec(`
		UPDATE reports 
		SET unit = ?, substations = ?, work_time_start = ?, work_time_end = ?,
			responsible_person = ?, work_content = ?, related_devices = ?,
			tools = ?, operation = ?, safety_measures = ?, remarks = ?,
			status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		report.Unit, string(substationsJSON), report.WorkTimeStart, report.WorkTimeEnd,
		report.ResponsiblePerson, report.WorkContent, string(devicesJSON),
		string(toolsJSON), report.Operation, report.SafetyMeasures,
		report.Remarks, report.Status, report.ID,
	)
	return err
}

func (r *ReportRepository) UpdateReportStatus(id int, status string) error {
	_, err := r.db.Exec(`
		UPDATE reports 
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		status, id,
	)
	return err
}

func (r *ReportRepository) ApproveReport(reportID int, isApproved bool, reviewComments string, reviewerID int, reviewerName string) error {
	status := "rejected"
	if isApproved {
		status = "approved"
	}

	_, err := r.db.Exec(`
		UPDATE reports 
		SET status = ?, review_comments = ?, reviewer_id = ?, reviewer_name = ?, reviewed_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		status, reviewComments, reviewerID, reviewerName, reportID,
	)
	return err
}

func (r *ReportRepository) RejectReport(id int, reviewComments string, reviewerID int, reviewerName string) error {
	_, err := r.db.Exec(`
		UPDATE reports 
		SET status = 'rejected', review_comments = ?, reviewer_id = ?, reviewer_name = ?, reviewed_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		reviewComments, reviewerID, reviewerName, id,
	)
	return err
}

func (r *ReportRepository) DeleteReport(reportID int) error {
	_, err := r.db.Exec("DELETE FROM reports WHERE id = ?", reportID)
	return err
}

func (r *ReportRepository) GetReportsByUser(userID int) ([]models.Report, error) {
	query := `
		SELECT id, report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			related_devices, tools, operation, safety_measures, remarks,
			status, review_comments, reviewer_id, reviewer_name, reviewed_at,
			created_at, updated_at
		FROM reports 
		WHERE reporter_id = ? 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.parseReports(rows)
}

func (r *ReportRepository) GetPendingReports() ([]models.Report, error) {
	query := `
		SELECT id, report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			related_devices, tools, operation, safety_measures, remarks,
			status, review_comments, reviewer_id, reviewer_name, reviewed_at,
			created_at, updated_at
		FROM reports 
		WHERE status = 'submitted'
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.parseReports(rows)
}

func (r *ReportRepository) GetAllReports() ([]models.Report, error) {
	query := `
		SELECT id, report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			related_devices, tools, operation, safety_measures, remarks,
			status, review_comments, reviewer_id, reviewer_name, reviewed_at,
			created_at, updated_at
		FROM reports 
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.parseReports(rows)
}

func (r *ReportRepository) parseReports(rows *sql.Rows) ([]models.Report, error) {
	var reports []models.Report

	for rows.Next() {
		var report models.Report
		var substationsJSON, devicesJSON, toolsJSON string
		var reviewedAt sql.NullTime

		err := rows.Scan(
			&report.ID, &report.ReportNo, &report.ReporterID, &report.ReporterName,
			&report.Unit, &substationsJSON, &report.WorkTimeStart, &report.WorkTimeEnd,
			&report.ResponsiblePerson, &report.WorkContent, &devicesJSON, &toolsJSON,
			&report.Operation, &report.SafetyMeasures, &report.Remarks, &report.Status,
			&report.ReviewComments, &report.ReviewerID, &report.ReviewerName,
			&reviewedAt, &report.CreatedAt, &report.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 解析JSON字段
		json.Unmarshal([]byte(substationsJSON), &report.Substations)
		json.Unmarshal([]byte(devicesJSON), &report.RelatedDevices)
		json.Unmarshal([]byte(toolsJSON), &report.Tools)

		if reviewedAt.Valid {
			report.ReviewedAt = reviewedAt
		}

		reports = append(reports, report)
	}

	return reports, nil
}

// 添加缺失的方法
func (r *ReportRepository) CreateReportLog(log *models.ReportLog) error {
	_, err := r.db.Exec(`
		INSERT INTO report_logs (report_id, action, operator_id, operator_name, comment, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		log.ReportID, log.Action, log.OperatorID, log.OperatorName, log.Comment, log.CreatedAt)
	return err
}

func (r *ReportRepository) GetReportsByUserWithFilter(userID int, status, startDate, endDate, reportNo, substation string) ([]models.Report, error) {
	query := `
		SELECT id, report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			related_devices, tools, operation, safety_measures, remarks,
			status, review_comments, reviewer_id, reviewer_name, reviewed_at,
			created_at, updated_at
		FROM reports 
		WHERE reporter_id = ?`

	args := []interface{}{userID}

	// 添加筛选条件
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if reportNo != "" {
		query += " AND report_no LIKE ?"
		args = append(args, "%"+reportNo+"%")
	}

	if startDate != "" {
		query += " AND DATE(work_time_start) >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND DATE(work_time_end) <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.parseReports(rows)
}

func (r *ReportRepository) GetAllReportsWithFilter(reporterName, reportNo, substation, startDate, endDate, status string) ([]models.Report, error) {
	query := `
		SELECT id, report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			related_devices, tools, operation, safety_measures, remarks,
			status, review_comments, reviewer_id, reviewer_name, reviewed_at,
			created_at, updated_at
		FROM reports 
		WHERE 1=1`

	var args []interface{}

	// 添加筛选条件
	if reporterName != "" {
		query += " AND reporter_name LIKE ?"
		args = append(args, "%"+reporterName+"%")
	}

	if reportNo != "" {
		query += " AND report_no LIKE ?"
		args = append(args, "%"+reportNo+"%")
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if startDate != "" {
		query += " AND DATE(work_time_start) >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND DATE(work_time_end) <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.parseReports(rows)
}

func (r *ReportRepository) GetReportLogs(reportID int) ([]models.ReportLog, error) {
	rows, err := r.db.Query(`
		SELECT id, report_id, action, operator_id, operator_name, comment, created_at
		FROM report_logs
		WHERE report_id = ?
		ORDER BY created_at DESC`,
		reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ReportLog
	for rows.Next() {
		var log models.ReportLog
		err := rows.Scan(&log.ID, &log.ReportID, &log.Action, &log.OperatorID,
			&log.OperatorName, &log.Comment, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *ReportRepository) GetReportStatistics(startDate, endDate string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'draft' THEN 1 ELSE 0 END) as draft_count,
			SUM(CASE WHEN status = 'submitted' THEN 1 ELSE 0 END) as submitted_count,
			SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END) as approved_count,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected_count
		FROM reports`

	var args []interface{}
	if startDate != "" && endDate != "" {
		query += " WHERE DATE(created_at) BETWEEN ? AND ?"
		args = append(args, startDate, endDate)
	}

	var total, draft, submitted, approved, rejected int
	err := r.db.QueryRow(query, args...).Scan(&total, &draft, &submitted, &approved, &rejected)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":         total,
		"draft":         draft,
		"submitted":     submitted,
		"approved":      approved,
		"rejected":      rejected,
		"approval_rate": fmt.Sprintf("%.1f%%", float64(approved)/float64(total)*100),
	}, nil
}

func (r *ReportRepository) GetDashboardStats() (map[string]interface{}, error) {
	// 获取今日报备数
	var todayCount int
	today := time.Now().Format("2006-01-02")
	r.db.QueryRow("SELECT COUNT(*) FROM reports WHERE DATE(created_at) = ?", today).Scan(&todayCount)

	// 获取待审核数
	var pendingCount int
	r.db.QueryRow("SELECT COUNT(*) FROM reports WHERE status = 'submitted'").Scan(&pendingCount)

	// 获取总报备数
	var totalCount int
	r.db.QueryRow("SELECT COUNT(*) FROM reports").Scan(&totalCount)

	// 获取最近7天数据
	var weeklyData []map[string]interface{}
	rows, err := r.db.Query(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM reports
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date string
			var count int
			rows.Scan(&date, &count)
			weeklyData = append(weeklyData, map[string]interface{}{
				"date":  date,
				"count": count,
			})
		}
	}

	return map[string]interface{}{
		"today_count":   todayCount,
		"pending_count": pendingCount,
		"total_count":   totalCount,
		"weekly_data":   weeklyData,
	}, nil
}

func (r *ReportRepository) GetRecentReports(limit int) ([]models.Report, error) {
	query := `
		SELECT id, report_no, reporter_id, reporter_name, unit, substations,
			work_time_start, work_time_end, responsible_person, work_content,
			status, created_at
		FROM reports 
		ORDER BY created_at DESC
		LIMIT ?`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var report models.Report
		var substationsJSON string
		err := rows.Scan(&report.ID, &report.ReportNo, &report.ReporterID,
			&report.ReporterName, &report.Unit, &substationsJSON,
			&report.WorkTimeStart, &report.WorkTimeEnd, &report.ResponsiblePerson,
			&report.WorkContent, &report.Status, &report.CreatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(substationsJSON), &report.Substations)
		reports = append(reports, report)
	}
	return reports, nil
}

func (r *ReportRepository) GetTemplateByID(id int) (*models.Template, error) {
	var template models.Template
	var substationsJSON, devicesJSON, toolsJSON string

	query := `
		SELECT id, name, unit, substations, work_content, related_devices, 
			tools, operation, safety_measures, created_by, created_at, updated_at
		FROM templates 
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&template.ID, &template.Name, &template.Unit, &substationsJSON,
		&template.WorkContent, &devicesJSON, &toolsJSON, &template.Operation,
		&template.SafetyMeasures, &template.CreatedBy, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// 解析JSON字段
	json.Unmarshal([]byte(substationsJSON), &template.Substations)
	json.Unmarshal([]byte(devicesJSON), &template.RelatedDevices)
	json.Unmarshal([]byte(toolsJSON), &template.Tools)

	return &template, nil
}

func (r *ReportRepository) UpdateUserProfile(user *models.User) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET real_name = ?, phone = ?, position = ?, unit = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		user.RealName, user.Phone, user.Position, user.Unit, user.ID)
	return err
}

func (r *ReportRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		SELECT id, username, password, real_name, phone, position, role, unit
		FROM users
		WHERE id = ?`, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.RealName,
		&user.Phone, &user.Position, &user.Role, &user.Unit)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
