package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Migration struct {
	Version int64
	Name    string
	UpSQL   string
	DownSQL string
}

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) CreateMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		version BIGINT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	return nil
}

func (m *Migrator) GetAppliedMigrations() (map[int64]bool, error) {
	applied := make(map[int64]bool)
	
	rows, err := m.db.Query("SELECT version FROM migrations ORDER BY version")
	if err != nil {
		return applied, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return applied, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}
	
	return applied, nil
}

func (m *Migrator) LoadMigrations(migrationsDir string) ([]Migration, error) {
	var migrations []Migration
	
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return migrations, fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}
		
		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Printf("Warning: invalid migration filename format: %s", file.Name())
			continue
		}
		
		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")
		
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return migrations, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}
		
		sqlContent := string(content)
		upSQL, downSQL := parseMigrationContent(sqlContent)
		
		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			UpSQL:   upSQL,
			DownSQL: downSQL,
		})
	}
	
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	
	return migrations, nil
}

func parseMigrationContent(content string) (string, string) {
	parts := strings.Split(content, "-- +migrate Down")
	upSQL := strings.TrimSpace(strings.Replace(parts[0], "-- +migrate Up", "", 1))
	
	downSQL := ""
	if len(parts) > 1 {
		downSQL = strings.TrimSpace(parts[1])
	}
	
	return upSQL, downSQL
}

func (m *Migrator) Up(migrationsDir string) error {
	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}
	
	migrations, err := m.LoadMigrations(migrationsDir)
	if err != nil {
		return err
	}
	
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}
	
	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}
		
		log.Printf("Applying migration %d: %s", migration.Version, migration.Name)
		
		if _, err := m.db.Exec(migration.UpSQL); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}
		
		if _, err := m.db.Exec("INSERT INTO migrations (version, name) VALUES (?, ?)", migration.Version, migration.Name); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}
		
		log.Printf("Migration %d applied successfully", migration.Version)
	}
	
	log.Println("All migrations applied successfully")
	return nil
}

func (m *Migrator) Down(migrationsDir string, steps int) error {
	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}
	
	migrations, err := m.LoadMigrations(migrationsDir)
	if err != nil {
		return err
	}
	
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}
	
	var toRollback []Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if applied[migrations[i].Version] {
			toRollback = append(toRollback, migrations[i])
			if len(toRollback) >= steps {
				break
			}
		}
	}
	
	for _, migration := range toRollback {
		if migration.DownSQL == "" {
			log.Printf("Warning: No down migration for %d: %s", migration.Version, migration.Name)
			continue
		}
		
		log.Printf("Rolling back migration %d: %s", migration.Version, migration.Name)
		
		if _, err := m.db.Exec(migration.DownSQL); err != nil {
			return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
		}
		
		if _, err := m.db.Exec("DELETE FROM migrations WHERE version = ?", migration.Version); err != nil {
			return fmt.Errorf("failed to remove migration record %d: %w", migration.Version, err)
		}
		
		log.Printf("Migration %d rolled back successfully", migration.Version)
	}
	
	log.Printf("Rolled back %d migrations successfully", len(toRollback))
	return nil
}

func (m *Migrator) Status(migrationsDir string) error {
	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}
	
	migrations, err := m.LoadMigrations(migrationsDir)
	if err != nil {
		return err
	}
	
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}
	
	fmt.Println("Migration Status:")
	fmt.Println("================")
	
	if len(migrations) == 0 {
		fmt.Println("No migrations found")
		return nil
	}
	
	for _, migration := range migrations {
		status := "Pending"
		if applied[migration.Version] {
			status = "Applied"
		}
		fmt.Printf("%d\t%s\t%s\n", migration.Version, status, migration.Name)
	}
	
	return nil
}

func CreateMigration(migrationsDir, name string) error {
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}
	
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, strings.ReplaceAll(name, " ", "_"))
	filepath := filepath.Join(migrationsDir, filename)
	
	content := fmt.Sprintf(`-- +migrate Up
-- Add your up migration here

-- +migrate Down
-- Add your down migration here
`)
	
	if err := ioutil.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	
	log.Printf("Migration created: %s", filepath)
	return nil
}