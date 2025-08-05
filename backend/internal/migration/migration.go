// Package migration handles running migrations at startup
package migration

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const MaxMigrations = 20

type MigrationResult struct {
	Name  string
	Count int
}
type Migration struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

type MigrationType string

const (
	MigrationTypeUp   MigrationType = "up"
	MigrationTypeDown MigrationType = "down"
)

type MigrationService struct {
	db *pgxpool.Pool
}

func NewMigration(db *pgxpool.Pool) *MigrationService {
	return &MigrationService{db: db}
}

func (m *MigrationService) Up() error {
	log.Println("Running migrations up")
	dbMigrations, err := m.GetDBMigrations()
	if err != nil {
		return err
	}

	localMigrations, err := m.GetLocalMigrations()
	if err != nil {
		return err
	}

	successFiles := make([]MigrationResult, 0, len(localMigrations))
	for _, migration := range localMigrations {
		if !slices.Contains(dbMigrations, migration) {
			count, err := m.runMigration(migration, MigrationTypeUp)
			if err != nil {
				return err
			}
			successFiles = append(successFiles, MigrationResult{
				Name:  migration,
				Count: count,
			})
		}
	}

	for _, migration := range successFiles {
		err = m.createMigration(migration.Name)
		if err != nil {
			return err
		}
	}

	for _, successFile := range successFiles {
		log.Printf("Migrations up '%s' completed successfully: %d queries applied\n", successFile.Name, successFile.Count)
	}

	return nil
}

func (m *MigrationService) Down() error {
	log.Println("Running migrations down")
	dbMigrations, err := m.GetDBMigrations()
	if err != nil {
		return err
	}

	localMigrations, err := m.GetLocalMigrations()
	if err != nil {
		return err
	}
	if len(dbMigrations) == 0 {
		log.Println("No migrations found in DB to roll down.")
		return nil
	}
	lastMigration := dbMigrations[len(dbMigrations)-1]

	successFiles := make([]MigrationResult, 0, 1)

	for _, migration := range localMigrations {
		if migration == lastMigration {
			count, err := m.runMigration(migration, MigrationTypeDown)
			if err != nil {
				return err
			}
			successFiles = append(successFiles, MigrationResult{
				Name:  migration,
				Count: count,
			})
		}
	}

	for _, migration := range successFiles {
		err = m.removeMigration(migration.Name)
		if err != nil {
			return err
		}
	}

	for _, successFile := range successFiles {
		log.Printf("Migrations down '%s' completed successfully: %d queries applied\n", successFile.Name, successFile.Count)
	}

	return nil
}

func (m *MigrationService) GetDBMigrations() ([]string, error) {
	query := `SELECT name FROM migrations`
	rows, err := m.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var migrations = make([]string, 0, MaxMigrations)

	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func (m *MigrationService) GetLocalMigrations() ([]string, error) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	migrations := make([]string, 0, len(files))
	for _, file := range files {
		migrations = append(migrations, file.Name())
	}

	return migrations, nil
}
func (m *MigrationService) runMigration(migration string, migrationType MigrationType) (int, error) {
	filePath := fmt.Sprintf("migrations/%s/%s.sql", migration, migrationType)
	fileExists, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("migration file %s not found", migration)
	}
	if fileExists.IsDir() {
		return 0, fmt.Errorf("migration file %s is a directory", migration)
	}

	query, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	queries := strings.Split(string(query), ";")
	tx, err := m.db.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	for _, query := range queries {
		if strings.TrimSpace(query) == "" {
			continue
		}

		_, err = tx.Exec(context.Background(), query)
		if err != nil {
			return 0, fmt.Errorf("failed to execute query: `%s` error: %w", query, err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return len(queries), nil
}
func (m *MigrationService) createMigration(migration string) error {
	query := `INSERT INTO migrations (name, created_at) VALUES ($1, $2)`
	_, err := m.db.Exec(context.Background(), query, migration, time.Now())
	if err != nil {
		return err
	}
	return nil
}
func (m *MigrationService) removeMigration(migration string) error {
	query := `DELETE FROM migrations WHERE name = $1`
	_, err := m.db.Exec(context.Background(), query, migration)
	if err != nil {
		return err
	}
	return nil
}
