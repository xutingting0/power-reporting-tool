// file: repositories/report_log_repo.go
package repositories

import (
	"backend/models"
	"database/sql"
)

type ReportLogRepository struct {
	db *sql.DB
}

func NewReportLogRepository(db *sql.DB) *ReportLogRepository {
	return &ReportLogRepository{db: db}
}

func (r *ReportLogRepository) CreateReportLog(log *models.ReportLog) error {
	_, err := r.db.Exec(`
		INSERT INTO report_logs (report_id, action, operator_id, operator_name, comment)
		VALUES (?, ?, ?, ?, ?)`,
		log.ReportID, log.Action, log.OperatorID, log.OperatorName, log.Comment,
	)
	return err
}

func (r *ReportLogRepository) GetReportLogs(reportID int) ([]models.ReportLog, error) {
	rows, err := r.db.Query(`
		SELECT id, report_id, action, operator_id, operator_name, comment, created_at
		FROM report_logs 
		WHERE report_id = ?
		ORDER BY created_at DESC`, reportID)
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
