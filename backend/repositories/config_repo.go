// file: repositories/config_repo.go
package repositories

import (
	"backend/models"
	"database/sql"
	"encoding/json"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) GetAllUnits() ([]models.ConfigUnit, error) {
	rows, err := r.db.Query(`
		SELECT id, name, code, description 
		FROM config_units 
		ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []models.ConfigUnit
	for rows.Next() {
		var unit models.ConfigUnit
		err := rows.Scan(&unit.ID, &unit.Name, &unit.Code, &unit.Description)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}

func (r *ConfigRepository) AddUnit(unit *models.ConfigUnit) error {
	_, err := r.db.Exec(`
		INSERT INTO config_units (name, code, description) 
		VALUES (?, ?, ?)`,
		unit.Name, unit.Code, unit.Description)
	return err
}

func (r *ConfigRepository) GetAllSubstations() ([]models.ConfigSubstation, error) {
	rows, err := r.db.Query(`
		SELECT id, name, voltage_level, region, description, status 
		FROM config_substations 
		WHERE status = 'active'
		ORDER BY voltage_level, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var substations []models.ConfigSubstation
	for rows.Next() {
		var s models.ConfigSubstation
		err := rows.Scan(&s.ID, &s.Name, &s.VoltageLevel, &s.Region, &s.Description, &s.Status)
		if err != nil {
			return nil, err
		}
		substations = append(substations, s)
	}
	return substations, nil
}

func (r *ConfigRepository) AddSubstation(substation *models.ConfigSubstation) error {
	_, err := r.db.Exec(`
		INSERT INTO config_substations (name, voltage_level, region, description, status) 
		VALUES (?, ?, ?, ?, ?)`,
		substation.Name, substation.VoltageLevel, substation.Region,
		substation.Description, substation.Status)
	return err
}

func (r *ConfigRepository) GetAllSafetyTemplates() ([]models.ConfigSafetyTemplate, error) {
	rows, err := r.db.Query(`
		SELECT id, name, content, is_required, template_type, sort_order 
		FROM config_safety_templates 
		ORDER BY sort_order, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []models.ConfigSafetyTemplate
	for rows.Next() {
		var t models.ConfigSafetyTemplate
		err := rows.Scan(&t.ID, &t.Name, &t.Content, &t.IsRequired, &t.Type, &t.SortOrder)
		if err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *ConfigRepository) GetRequiredSafetyMeasures() (string, error) {
	rows, err := r.db.Query(`
		SELECT content 
		FROM config_safety_templates 
		WHERE is_required = 1 
		ORDER BY sort_order`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var measures []string
	for rows.Next() {
		var content string
		err := rows.Scan(&content)
		if err != nil {
			return "", err
		}
		measures = append(measures, content)
	}

	result, _ := json.Marshal(measures)
	return string(result), nil
}

func (r *ConfigRepository) GetAllDeviceTypes() ([]models.ConfigDeviceType, error) {
	rows, err := r.db.Query(`
		SELECT id, name, category, description 
		FROM config_device_types 
		ORDER BY category, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.ConfigDeviceType
	for rows.Next() {
		var t models.ConfigDeviceType
		err := rows.Scan(&t.ID, &t.Name, &t.Category, &t.Description)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

// 添加缺失的方法
func (r *ConfigRepository) GetSubstationsByVoltageLevel() (map[string][]models.ConfigSubstation, error) {
	rows, err := r.db.Query(`
		SELECT id, name, voltage_level, region, description, status 
		FROM config_substations 
		WHERE status = 'active'
		ORDER BY voltage_level, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouped := make(map[string][]models.ConfigSubstation)
	for rows.Next() {
		var s models.ConfigSubstation
		err := rows.Scan(&s.ID, &s.Name, &s.VoltageLevel, &s.Region, &s.Description, &s.Status)
		if err != nil {
			return nil, err
		}
		grouped[s.VoltageLevel] = append(grouped[s.VoltageLevel], s)
	}
	return grouped, nil
}

func (r *ConfigRepository) DeleteUnit(id int) error {
	_, err := r.db.Exec("DELETE FROM config_units WHERE id = ?", id)
	return err
}

func (r *ConfigRepository) DeleteSubstation(id int) error {
	_, err := r.db.Exec("DELETE FROM config_substations WHERE id = ?", id)
	return err
}

func (r *ConfigRepository) UpdateUnit(unit *models.ConfigUnit) error {
	_, err := r.db.Exec(`
		UPDATE config_units 
		SET name = ?, code = ?, description = ?
		WHERE id = ?`,
		unit.Name, unit.Code, unit.Description, unit.ID)
	return err
}

func (r *ConfigRepository) UpdateSubstation(substation *models.ConfigSubstation) error {
	_, err := r.db.Exec(`
		UPDATE config_substations 
		SET name = ?, voltage_level = ?, region = ?, description = ?, status = ?
		WHERE id = ?`,
		substation.Name, substation.VoltageLevel, substation.Region,
		substation.Description, substation.Status, substation.ID)
	return err
}
