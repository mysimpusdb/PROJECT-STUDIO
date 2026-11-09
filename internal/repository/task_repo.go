package repository

import (
	"database/sql"
	"strings"
	"time"

	"project-studio/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// taskColumns is the canonical column list used by all SELECTs so that
// scanTask stays in sync with the queries.
const taskColumns = `id, project_id, parent_id, name, description, status, priority, start_date, end_date,
	duration, estimated_hours, progress, assignee_id, sort_order, dependencies, is_milestone,
	labels, checklists, estimated_cost, actual_cost, created_at, updated_at`

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func formatDate(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format("2006-01-02")
}

// parseDate accepts both "2006-01-02" and full datetime strings.
func parseDate(s sql.NullString) *time.Time {
	if !s.Valid || len(s.String) < 10 {
		return nil
	}
	if t, err := time.Parse("2006-01-02", s.String[:10]); err == nil {
		return &t
	}
	return nil
}

func parseDateTime(s sql.NullString) time.Time {
	if !s.Valid {
		return time.Time{}
	}
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02T15:04:05Z"} {
		if t, err := time.Parse(layout, s.String); err == nil {
			return t
		}
	}
	return time.Time{}
}

func scanTask(row rowScanner) (*models.Task, error) {
	var t models.Task
	var parentID, description, startDate, endDate, assigneeID, dependencies sql.NullString
	var labels, checklists, createdAt, updatedAt sql.NullString
	var estimatedHours sql.NullInt64
	var estimatedCost, actualCost sql.NullFloat64
	var isMilestone sql.NullBool

	err := row.Scan(
		&t.ID,
		&t.ProjectID,
		&parentID,
		&t.Name,
		&description,
		&t.Status,
		&t.Priority,
		&startDate,
		&endDate,
		&t.Duration,
		&estimatedHours,
		&t.Progress,
		&assigneeID,
		&t.SortOrder,
		&dependencies,
		&isMilestone,
		&labels,
		&checklists,
		&estimatedCost,
		&actualCost,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	if parentID.Valid && parentID.String != "" {
		t.ParentID = &parentID.String
	}
	if assigneeID.Valid && assigneeID.String != "" {
		t.AssigneeID = &assigneeID.String
	}
	t.Description = description.String
	t.Dependencies = dependencies.String
	t.Labels = labels.String
	t.Checklists = checklists.String
	t.EstimatedHours = int(estimatedHours.Int64)
	t.EstimatedCost = estimatedCost.Float64
	t.ActualCost = actualCost.Float64
	t.IsMilestone = isMilestone.Bool
	t.StartDate = parseDate(startDate)
	t.EndDate = parseDate(endDate)
	t.CreatedAt = parseDateTime(createdAt)
	t.UpdatedAt = parseDateTime(updatedAt)

	return &t, nil
}

// nullIfEmpty stores empty optional IDs as NULL so foreign keys stay valid.
func nullIfEmpty(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.createWith(r.db, task)
}

// CreateTx inserts a task inside an existing transaction.
func (r *TaskRepository) CreateTx(tx *sql.Tx, task *models.Task) error {
	return r.createWith(tx, task)
}

type execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (r *TaskRepository) createWith(ex execer, task *models.Task) error {
	query := `
		INSERT INTO tasks (id, project_id, parent_id, name, description, status, priority, start_date, end_date,
			duration, estimated_hours, progress, assignee_id, sort_order, dependencies, is_milestone,
			labels, checklists, estimated_cost, actual_cost)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := ex.Exec(query,
		task.ID,
		task.ProjectID,
		nullIfEmpty(task.ParentID),
		task.Name,
		task.Description,
		task.Status,
		task.Priority,
		formatDate(task.StartDate),
		formatDate(task.EndDate),
		task.Duration,
		task.EstimatedHours,
		task.Progress,
		nullIfEmpty(task.AssigneeID),
		task.SortOrder,
		task.Dependencies,
		task.IsMilestone,
		task.Labels,
		task.Checklists,
		task.EstimatedCost,
		task.ActualCost,
	)
	return err
}

func (r *TaskRepository) GetByID(id string) (*models.Task, error) {
	row := r.db.QueryRow(`SELECT `+taskColumns+` FROM tasks WHERE id = ?`, id)
	t, err := scanTask(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (r *TaskRepository) GetByProjectID(projectID string) ([]models.Task, error) {
	rows, err := r.db.Query(`
		SELECT `+taskColumns+`
		FROM tasks
		WHERE project_id = ?
		ORDER BY sort_order ASC, created_at ASC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) GetByAssigneeID(assigneeID string) ([]models.Task, error) {
	rows, err := r.db.Query(`
		SELECT `+taskColumns+`
		FROM tasks
		WHERE assignee_id = ?
		ORDER BY end_date ASC, created_at DESC
	`, assigneeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) Update(task *models.Task) error {
	query := `
		UPDATE tasks
		SET name = ?, description = ?, status = ?, priority = ?, start_date = ?, end_date = ?, duration = ?,
			estimated_hours = ?, progress = ?, assignee_id = ?, sort_order = ?, dependencies = ?, is_milestone = ?,
			labels = ?, checklists = ?, estimated_cost = ?, actual_cost = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		task.Name,
		task.Description,
		task.Status,
		task.Priority,
		formatDate(task.StartDate),
		formatDate(task.EndDate),
		task.Duration,
		task.EstimatedHours,
		task.Progress,
		nullIfEmpty(task.AssigneeID),
		task.SortOrder,
		task.Dependencies,
		task.IsMilestone,
		task.Labels,
		task.Checklists,
		task.EstimatedCost,
		task.ActualCost,
		task.ID,
	)
	return err
}

// Delete removes a task together with all of its subtasks and dependent rows
// (comments, attachments, time logs, assignees, dependencies) in a single
// transaction. It returns the file paths of attachments that were removed so
// the caller can delete the physical files.
func (r *TaskRepository) Delete(id string) ([]string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var projectID string
	if err := tx.QueryRow(`SELECT project_id FROM tasks WHERE id = ?`, id).Scan(&projectID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Collect the task and all its descendants.
	rows, err := tx.Query(`
		WITH RECURSIVE subtree(id) AS (
			SELECT id FROM tasks WHERE id = ?
			UNION ALL
			SELECT t.id FROM tasks t JOIN subtree s ON t.parent_id = s.id
		)
		SELECT id FROM subtree
	`, id)
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var tid string
		if err := rows.Scan(&tid); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, tid)
	}
	rows.Close()

	files, err := deleteTaskRows(tx, ids)
	if err != nil {
		return nil, err
	}

	// Remove references to deleted tasks from the remaining tasks' dependency lists.
	if err := pruneDependencies(tx, projectID, ids); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return files, nil
}

// deleteTaskRows deletes the given tasks and every row that references them.
func deleteTaskRows(tx *sql.Tx, ids []string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}

	var files []string
	rows, err := tx.Query(`SELECT file_path FROM task_attachments WHERE task_id IN (`+placeholders+`)`, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err == nil {
			files = append(files, p)
		}
	}
	rows.Close()

	stmts := []string{
		`DELETE FROM task_comments WHERE task_id IN (` + placeholders + `)`,
		`DELETE FROM task_attachments WHERE task_id IN (` + placeholders + `)`,
		`DELETE FROM time_logs WHERE task_id IN (` + placeholders + `)`,
		`DELETE FROM task_assignees WHERE task_id IN (` + placeholders + `)`,
		`DELETE FROM task_dependencies WHERE predecessor_task_id IN (` + placeholders + `) OR successor_task_id IN (` + placeholders + `)`,
		`DELETE FROM tasks WHERE id IN (` + placeholders + `)`,
	}
	for _, q := range stmts {
		a := args
		if strings.Count(q, "IN (") == 2 {
			a = append(append([]interface{}{}, args...), args...)
		}
		if _, err := tx.Exec(q, a...); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func pruneDependencies(tx *sql.Tx, projectID string, removed []string) error {
	removedSet := make(map[string]bool, len(removed))
	for _, id := range removed {
		removedSet[id] = true
	}

	rows, err := tx.Query(`SELECT id, dependencies FROM tasks WHERE project_id = ? AND dependencies IS NOT NULL AND dependencies != ''`, projectID)
	if err != nil {
		return err
	}
	type upd struct{ id, deps string }
	var updates []upd
	for rows.Next() {
		var id, deps string
		if err := rows.Scan(&id, &deps); err != nil {
			rows.Close()
			return err
		}
		var kept []string
		changed := false
		for _, d := range strings.Split(deps, ",") {
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			if removedSet[d] {
				changed = true
				continue
			}
			kept = append(kept, d)
		}
		if changed {
			updates = append(updates, upd{id, strings.Join(kept, ", ")})
		}
	}
	rows.Close()

	for _, u := range updates {
		if _, err := tx.Exec(`UPDATE tasks SET dependencies = ? WHERE id = ?`, u.deps, u.id); err != nil {
			return err
		}
	}
	return nil
}

// DeleteByProjectTx removes every task of a project (and dependent rows)
// inside an existing transaction. Returns attachment file paths.
func (r *TaskRepository) DeleteByProjectTx(tx *sql.Tx, projectID string) ([]string, error) {
	rows, err := tx.Query(`SELECT id FROM tasks WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	return deleteTaskRows(tx, ids)
}
