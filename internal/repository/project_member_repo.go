package repository

import (
	"database/sql"

	"project-studio/internal/models"
)

type ProjectMemberRepository struct {
	db *sql.DB
}

func NewProjectMemberRepository(db *sql.DB) *ProjectMemberRepository {
	return &ProjectMemberRepository{db: db}
}

// GetByProjectID returns all members for a project including user details
func (r *ProjectMemberRepository) GetByProjectID(projectID string) ([]models.ProjectMember, error) {
	query := `
		SELECT pm.project_id, pm.user_id, pm.role, u.email, u.full_name, u.avatar_url
		FROM project_members pm
		JOIN users u ON pm.user_id = u.id
		WHERE pm.project_id = ?
		ORDER BY u.full_name ASC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMember
	for rows.Next() {
		var m models.ProjectMember
		var avatar sql.NullString
		
		err := rows.Scan(
			&m.ProjectID,
			&m.UserID,
			&m.Role,
			&m.UserEmail,
			&m.UserFullName,
			&avatar,
		)
		if err != nil {
			return nil, err
		}
		
		if avatar.Valid {
			m.UserAvatar = &avatar.String
		}
		
		members = append(members, m)
	}
	
	return members, nil
}

// AddMember adds a user to a project with a specific role
func (r *ProjectMemberRepository) AddMember(projectID string, userID string, role string) error {
	query := `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES (?, ?, ?)
	`
	_, err := r.db.Exec(query, projectID, userID, role)
	return err
}

// UpdateRole updates a member's role in a project
func (r *ProjectMemberRepository) UpdateRole(projectID string, userID string, role string) error {
	query := `
		UPDATE project_members
		SET role = ?
		WHERE project_id = ? AND user_id = ?
	`
	_, err := r.db.Exec(query, role, projectID, userID)
	return err
}

// RemoveMember removes a user from a project
func (r *ProjectMemberRepository) RemoveMember(projectID string, userID string) error {
	query := `
		DELETE FROM project_members
		WHERE project_id = ? AND user_id = ?
	`
	_, err := r.db.Exec(query, projectID, userID)
	return err
}

// GetWorkspaceUsers returns all users in a workspace 
func (r *ProjectMemberRepository) GetWorkspaceUsers(workspaceID string) ([]models.User, error) {
	query := `
		SELECT u.id, u.email, u.full_name, u.avatar_url
		FROM workspace_members wm
		JOIN users u ON wm.user_id = u.id
		WHERE wm.workspace_id = ?
		ORDER BY u.full_name ASC
	`
	rows, err := r.db.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var avatar sql.NullString
		
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FullName,
			&avatar,
		)
		if err != nil {
			return nil, err
		}
		
		if avatar.Valid {
			u.AvatarURL = &avatar.String
		}
		
		users = append(users, u)
	}
	
	return users, nil
}
