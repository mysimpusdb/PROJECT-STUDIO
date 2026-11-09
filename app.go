package main

import (
	"context"
	"fmt"
	"log"
	"project-studio/internal/db"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	dbPath, err := db.DefaultDBPath()
	if err != nil {
		log.Fatalf("Error resolving database path: %v", err)
	}
	log.Printf("Using database: %s", dbPath)

	err = db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	
	err = db.SeedDefaultAdmin()
	if err != nil {
		log.Printf("Error seeding default admin: %v", err)
	}
}

// shutdown is called when the app closes
func (a *App) shutdown(ctx context.Context) {
	db.Close()
	log.Println("Database connection closed.")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
