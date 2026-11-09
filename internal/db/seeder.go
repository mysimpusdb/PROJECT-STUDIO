package db

import (
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SeedDefaultAdmin creates a default admin user if no users exist
func SeedDefaultAdmin() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Users already exist, no need to seed
	}

	adminEmail := "admin@dbstudio.local"
	adminPassword := "admin123"
	adminName := "System Admin"

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminID := uuid.New().String()
	_, err = DB.Exec(`
		INSERT INTO users (id, email, password_hash, full_name) 
		VALUES (?, ?, ?, ?)
	`, adminID, adminEmail, string(hash), adminName)

	if err != nil {
		return err
	}

	// Create a default workspace for the admin
	workspaceID := uuid.New().String()
	_, err = DB.Exec(`
		INSERT INTO workspaces (id, name, owner_id) 
		VALUES (?, ?, ?)
	`, workspaceID, "Default Workspace", adminID)

	if err != nil {
		return err
	}

	// Add admin to workspace members
	_, err = DB.Exec(`
		INSERT INTO workspace_members (workspace_id, user_id, role) 
		VALUES (?, ?, 'admin')
	`, workspaceID, adminID)

	if err != nil {
		return err
	}

	log.Printf("Successfully seeded default admin: %s\n", adminEmail)
	return nil
}
