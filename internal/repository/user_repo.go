package repository

import (
	"database/sql"
	"errors"
	"time"

	"project-studio/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	tz := user.Timezone
	if tz == "" {
		tz = "UTC"
	}
	_, err := r.db.Exec(`
		INSERT INTO users (id, email, password_hash, full_name, avatar_url, timezone)
		VALUES (?, ?, ?, ?, ?, ?)
	`, user.ID, user.Email, user.PasswordHash, user.FullName, user.AvatarURL, tz)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, timezone, created_at, updated_at
		FROM users
		WHERE email = ?
	`
	row := r.db.QueryRow(query, email)

	var user models.User
	var avatarURL sql.NullString
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&avatarURL,
		&user.Timezone,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}

	// SQLite returns datetime as string, parse it
	layout := "2006-01-02 15:04:05" // SQLite standard datetime format
	
	if t, err := time.Parse(layout, createdAtStr); err == nil {
		user.CreatedAt = t
	}
	if t, err := time.Parse(layout, updatedAtStr); err == nil {
		user.UpdatedAt = t
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, timezone, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var user models.User
	var avatarURL sql.NullString
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&avatarURL,
		&user.Timezone,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}

	layout := "2006-01-02 15:04:05"
	if t, err := time.Parse(layout, createdAtStr); err == nil {
		user.CreatedAt = t
	}
	if t, err := time.Parse(layout, updatedAtStr); err == nil {
		user.UpdatedAt = t
	}

	return &user, nil
}

func (r *UserRepository) UpdateProfile(id, fullName, timezone, notifPrefs string, avatarURL *string) error {
	query := `
		UPDATE users 
		SET full_name = ?, timezone = ?, notification_preferences = ?, avatar_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query, fullName, timezone, notifPrefs, avatarURL, id)
	return err
}

func (r *UserRepository) UpdatePassword(id, passwordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query, passwordHash, id)
	return err
}
