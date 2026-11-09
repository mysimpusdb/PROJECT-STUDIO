package repository

import (
	"database/sql"

	"project-studio/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

const projectColumns = `id, workspace_id, name, description, color, status, start_date, end_date,
	is_public, share_token, total_budget, currency, created_by, created_at, updated_at`

func scanProject(row rowScanner) (*models.Project, error) {
	var p models.Project
	var description, color, startDate, endDate, shareToken, currency sql.NullString
	var createdAt, updatedAt sql.NullString
	var isPublic sql.NullBool
	var totalBudget sql.NullFloat64

	err := row.Scan(
		&p.ID,
		&p.WorkspaceID,
		&p.Name,
		&description,
		&color,
		&p.Status,
		&startDate,
		&endDate,
		&isPublic,
		&shareToken,
		&totalBudget,
		&currency,
		&p.CreatedBy,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Description = description.String
	p.Color = color.String
	p.IsPublic = isPublic.Bool
	if shareToken.Valid && shareToken.String != "" {
		p.ShareToken = &shareToken.String
	}
	p.TotalBudget = totalBudget.Float64
	p.Currency = currency.String
	if p.Currency == "" {
		p.Currency = "IDR"
	}
	p.StartDate = parseDate(startDate)
	p.EndDate = parseDate(endDate)
	p.CreatedAt = parseDateTime(createdAt)
	p.UpdatedAt = parseDateTime(updatedAt)

	return &p, nil
}

func (r *ProjectRepository) Create(project *models.Project) error {
	if project.Currency == "" {
		project.Currency = "IDR"
	}
	query := `
		INSERT INTO projects (id, workspace_id, name, description, color, status, start_date, end_date,
			is_public, share_token, total_budget, currency, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		project.ID,
		project.WorkspaceID,
		project.Name,
		project.Description,
		project.Color,
		project.Status,
		formatDate(project.StartDate),
		formatDate(project.EndDate),
		project.IsPublic,
		nullIfEmpty(project.ShareToken),
		project.TotalBudget,
		project.Currency,
		project.CreatedBy,
	)
	return err
}

func (r *ProjectRepository) GetByWorkspaceID(workspaceID string) ([]models.Project, error) {
	rows, err := r.db.Query(`
		SELECT `+projectColumns+`
		FROM projects
		WHERE workspace_id = ?
		ORDER BY created_at DESC
	`, workspaceID)
	return collectProjects(rows, err)
}

// GetAccessibleByWorkspaceID returns only the projects in the workspace that
// the user created or is a member of.
func (r *ProjectRepository) GetAccessibleByWorkspaceID(workspaceID, userID string) ([]models.Project, error) {
	rows, err := r.db.Query(`
		SELECT `+projectColumns+`
		FROM projects
		WHERE workspace_id = ?
		  AND (created_by = ? OR id IN (SELECT project_id FROM project_members WHERE user_id = ?))
		ORDER BY created_at DESC
	`, workspaceID, userID, userID)
	return collectProjects(rows, err)
}

func collectProjects(rows *sql.Rows, err error) ([]models.Project, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) GetByID(id string) (*models.Project, error) {
	row := r.db.QueryRow(`SELECT `+projectColumns+` FROM projects WHERE id = ?`, id)
	p, err := scanProject(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *ProjectRepository) GetByShareToken(token string) (*models.Project, error) {
	row := r.db.QueryRow(`SELECT `+projectColumns+` FROM projects WHERE share_token = ? AND is_public = 1`, token)
	p, err := scanProject(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *ProjectRepository) Update(project *models.Project) error {
	if project.Currency == "" {
		project.Currency = "IDR"
	}
	query := `
		UPDATE projects
		SET name = ?, description = ?, color = ?, status = ?, start_date = ?, end_date = ?,
			is_public = ?, share_token = ?, total_budget = ?, currency = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		project.Name,
		project.Description,
		project.Color,
		project.Status,
		formatDate(project.StartDate),
		formatDate(project.EndDate),
		project.IsPublic,
		nullIfEmpty(project.ShareToken),
		project.TotalBudget,
		project.Currency,
		project.ID,
	)
	return err
}

// Delete permanently removes a project together with all its tasks, members
// and dependent rows in a single transaction. It returns attachment file paths
// so the caller can remove the physical files.
func (r *ProjectRepository) Delete(id string) ([]string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	files, err := NewTaskRepository(r.db).DeleteByProjectTx(tx, id)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`DELETE FROM project_members WHERE project_id = ?`, id); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`DELETE FROM projects WHERE id = ?`, id); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return files, nil
}
