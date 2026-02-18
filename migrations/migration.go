package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (

	//go:embed *.up.sql
	engineMigrationFS embed.FS
	appMigrationFS    fs.FS
	DB                *sqlx.DB
)

func NewMigration(db *sqlx.DB, fs fs.FS) {
	appMigrationFS = fs
	DB = db
}

func parseMigrationId(filename string) (string, error) {
	if !strings.HasSuffix(filename, ".up.sql") {
		return "", fmt.Errorf("invalid migration filename %q: must end with .up.sql", filename)
	}

	id := strings.TrimSuffix(filename, ".up.sql")

	s := id[:4]

	for _, v := range s {
		if v < '0' || v > '9' {
			return "", fmt.Errorf("invalid migration filename %q: must start with 4-digit prefix", filename)
		}
	}

	if id[4] != '_' {
		return "", fmt.Errorf("invalid migration filename %q: digit prefix must be followed by underscore", filename)
	}

	return id, nil
}

type migrationEntry struct {
	id       string
	filename string
	fsys     fs.FS
}

func getMigrationsEntry(fsys fs.FS) ([]*migrationEntry, error) {
	if fsys == nil {
		return nil, nil
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read mirations directory: %w", err)
	}

	results := []*migrationEntry{}

	for _, entrie := range entries {
		if entrie.IsDir() {
			continue
		}

		if !strings.HasSuffix(entrie.Name(), ".up.sql") {
			continue
		}

		id, err := parseMigrationId(entrie.Name())
		if err != nil {
			return nil, err
		}

		me := migrationEntry{
			id:       id,
			filename: entrie.Name(),
			fsys:     fsys,
		}

		results = append(results, &me)

	}

	return results, nil
}

func checkIfTableExists(tx *sql.Tx) (bool, error) {
	const query = `SELECT count(*) FROM information_schema.tables
	WHERE table_schema = 'public' AND table_name = 'schema_migrations'`

	var count int
	if err := tx.QueryRow(query).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check if schema_migrations exists: %w", err)
	}

	return count > 0, nil
}

func createMigrationsTable(tx *sql.Tx) error {
	const query = `CREATE TABLE IF NOT EXISTS schema_migrations
	(id TEXT PRIMARY KEY, applied_at TEXT DEFAULT CURRENT_TIMESTAMP)`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to create schema_migrations: %w", err)
	}

	return nil
}

func getAppliedMigrations(tx *sql.Tx) (map[string]bool, error) {
	const query = `SELECT id FROM schema_migrations`

	rows, err := tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}

	defer rows.Close()

	applied := make(map[string]bool)

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan migration id: %w", err)
		}

		applied[id] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating applied migrations: %w", err)
	}

	return applied, nil
}

func recordMigration(id string, tx *sql.Tx) error {
	const query = `INSERT INTO schema_migrations (id) VALUES ($1);`
	if _, err := tx.Exec(query, id); err != nil {
		return fmt.Errorf("failed to record migration %q: %w", id, err)
	}

	return nil
}

func Run() error {
	if DB == nil {
		return errors.New("DB *sqlx.DB is nil")
	}

	appMigrations, err := getMigrationsEntry(appMigrationFS)
	if err != nil {
		return fmt.Errorf("failed to collect application migrations: %w", err)
	}

	engineMigrations, err := getMigrationsEntry(engineMigrationFS)
	if err != nil {
		return fmt.Errorf("failed to collect engine migrations: %w", err)
	}

	allMigrations := append(appMigrations, engineMigrations...)

	seen := make(map[string]string)

	for _, m := range allMigrations {
		if existing, ok := seen[m.id]; ok {
			return fmt.Errorf("duplicate migration id %q: found in %q and %q", m.id, existing, m.filename)
		}

		seen[m.id] = m.filename
	}

	sort.Slice(allMigrations, func(i, j int) bool {
		return allMigrations[i].id < allMigrations[j].id
	})

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin trassaction: %w", err)
	}

	defer func() {
		if tx != nil {
			rberr := tx.Rollback()
			if rberr != nil {
				log.Printf("failed to rollback transcation: %v", rberr)
			}
		}
	}()

	exists, err := checkIfTableExists(tx)

	if !exists {
		err = createMigrationsTable(tx)
		if err != nil {
			return fmt.Errorf("failed to ensure schema_migrations table exists: %w", err)
		}
	}

	applied, err := getAppliedMigrations(tx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedCount := 0

	for _, m := range allMigrations {
		if applied[m.id] {
			continue
		}

		content, err := fs.ReadFile(m.fsys, m.filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %q: %w", m.filename, err)
		}

		log.Printf("applying migrations: %s", m.id)

		if _, err := tx.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to apply migration %q: %w", m.id, err)
		}

		if err := recordMigration(m.id, tx); err != nil {
			return err
		}

		appliedCount++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	tx = nil

	if appliedCount == 0 {
		log.Printf("no new migrations to apply\n")
	} else {
		log.Printf("applied %d migration(s)", appliedCount)
	}

	return nil
}
