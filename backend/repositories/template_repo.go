// file: repositories/template_repo.go
package repositories

import (
	"backend/models"
	"database/sql"
	"encoding/json"
)

type TemplateRepository struct {
	db *sql.DB
}

func NewTemplateRepository(db *sql.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) GetAllTemplates() ([]models.Template, error) {
	rows, err := r.db.Query(`
		SELECT id, name, unit, substations, work_content, related_devices, tools, 
			operation, safety_measures, created_by, created_at, updated_at
		FROM templates 
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var template models.Template
		var substationsJSON, devicesJSON, toolsJSON string

		err := rows.Scan(
			&template.ID, &template.Name, &template.Unit, &substationsJSON,
			&template.WorkContent, &devicesJSON, &toolsJSON, &template.Operation,
			&template.SafetyMeasures, &template.CreatedBy, &template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 解析JSON字段
		json.Unmarshal([]byte(substationsJSON), &template.Substations)
		json.Unmarshal([]byte(devicesJSON), &template.RelatedDevices)
		json.Unmarshal([]byte(toolsJSON), &template.Tools)

		templates = append(templates, template)
	}
	return templates, nil
}

func (r *TemplateRepository) GetTemplateByID(id int) (*models.Template, error) {
	var template models.Template
	var substationsJSON, devicesJSON, toolsJSON string

	query := `
		SELECT id, name, unit, substations, work_content, related_devices, tools, 
			operation, safety_measures, created_by, created_at, updated_at
		FROM templates 
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&template.ID, &template.Name, &template.Unit, &substationsJSON,
		&template.WorkContent, &devicesJSON, &toolsJSON, &template.Operation,
		&template.SafetyMeasures, &template.CreatedBy, &template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 解析JSON字段
	json.Unmarshal([]byte(substationsJSON), &template.Substations)
	json.Unmarshal([]byte(devicesJSON), &template.RelatedDevices)
	json.Unmarshal([]byte(toolsJSON), &template.Tools)

	return &template, nil
}

func (r *TemplateRepository) CreateTemplate(template *models.Template) error {
	substationsJSON, _ := json.Marshal(template.Substations)
	devicesJSON, _ := json.Marshal(template.RelatedDevices)
	toolsJSON, _ := json.Marshal(template.Tools)

	result, err := r.db.Exec(`
		INSERT INTO templates (name, unit, substations, work_content, related_devices, tools, operation, safety_measures, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		template.Name, template.Unit, string(substationsJSON), template.WorkContent,
		string(devicesJSON), string(toolsJSON), template.Operation,
		template.SafetyMeasures, template.CreatedBy,
	)

	if err != nil {
		return err
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err == nil {
		template.ID = int(id)
	}

	return nil
}

func (r *TemplateRepository) UpdateTemplate(template *models.Template) error {
	substationsJSON, _ := json.Marshal(template.Substations)
	devicesJSON, _ := json.Marshal(template.RelatedDevices)
	toolsJSON, _ := json.Marshal(template.Tools)

	_, err := r.db.Exec(`
		UPDATE templates 
		SET name = ?, unit = ?, substations = ?, work_content = ?, 
			related_devices = ?, tools = ?, operation = ?, safety_measures = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		template.Name, template.Unit, string(substationsJSON), template.WorkContent,
		string(devicesJSON), string(toolsJSON), template.Operation,
		template.SafetyMeasures, template.ID,
	)
	return err
}

func (r *TemplateRepository) DeleteTemplate(id int) error {
	_, err := r.db.Exec("DELETE FROM templates WHERE id = ?", id)
	return err
}
