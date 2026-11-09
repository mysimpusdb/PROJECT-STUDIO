package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// migrationSQL is embedded into the binary so production builds (wails build)
// do not depend on the source tree being present next to the executable.
//
//go:embed migrations.sql
var migrationSQL string

const appDirName = "ProjectStudio"
const dbFileName = "project-studio.db"

// DefaultDBPath returns the per-user location of the application database,
// e.g. %AppData%\ProjectStudio\project-studio.db on Windows.
// If a legacy database exists in the current working directory and no database
// exists yet in the new location, it is copied over once so no data is lost.
func DefaultDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve user config dir: %w", err)
	}

	dir := filepath.Join(configDir, appDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create app data dir: %w", err)
	}

	target := filepath.Join(dir, dbFileName)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		if info, err := os.Stat(dbFileName); err == nil && info.Size() > 0 {
			if err := copyFile(dbFileName, target); err != nil {
				log.Printf("Warning: failed to migrate legacy database: %v", err)
			} else {
				log.Printf("Legacy database copied to %s", target)
			}
		}
	}

	return target, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(dst)
		return err
	}
	return out.Close()
}

// InitDB initializes the database connection and runs migrations
func InitDB(dbPath string) error {
	dsn := dbPath
	// Enable FK enforcement, wait on locks instead of failing immediately,
	// and use WAL for better concurrent read/write behaviour.
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	dsn += sep + "_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	if dbPath != ":memory:" {
		dsn += "&_pragma=journal_mode(WAL)"
	}

	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if dbPath == ":memory:" {
		// Each new connection to :memory: is a separate empty database.
		DB.SetMaxOpenConns(1)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully.")

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

var alterAddColumnRe = regexp.MustCompile(`(?is)ALTER\s+TABLE\s+(\w+)\s+ADD\s+COLUMN\s+(\w+)`)

func columnExists(table, column string) (bool, error) {
	var count int
	err := DB.QueryRow(
		`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?`, table, column,
	).Scan(&count)
	return count > 0, err
}

func runMigrations() error {
	statements := strings.Split(migrationSQL, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || isOnlyComments(stmt) {
			continue
		}

		// ALTER TABLE ... ADD COLUMN is not idempotent in SQLite, so skip it
		// when the column already exists instead of swallowing every error.
		if m := alterAddColumnRe.FindStringSubmatch(stmt); m != nil {
			exists, err := columnExists(m[1], m[2])
			if err != nil {
				return fmt.Errorf("checking column %s.%s: %w", m[1], m[2], err)
			}
			if exists {
				continue
			}
		}

		if _, err := DB.Exec(stmt); err != nil {
			return fmt.Errorf("migration statement failed: %w\n%s", err, stmt)
		}
	}

	log.Println("Database migrations applied successfully.")
	return nil
}

func isOnlyComments(stmt string) bool {
	for _, line := range strings.Split(stmt, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "--") {
			return false
		}
	}
	return true
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}
