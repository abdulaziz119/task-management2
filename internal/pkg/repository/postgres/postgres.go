package postgres

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"io/ioutil"
	"log"
	"strings"
	"task-management/internal/pkg/config"
)

func NewPostgres() *bun.DB {
	dsn := "postgres://" + config.GetConf().DBUsername + ":" + config.GetConf().DBPassword + "@" +
		config.GetConf().DBHost + ":" + config.GetConf().DBPort + "/" + config.GetConf().DBName +
		"?sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := runMigrations(db); err != nil {
		log.Printf("Error running migrations: %v", err)
	}

	return db
}

func runMigrations(db *bun.DB) error {
	migrationFiles := []string{
		"internal/pkg/script/migrations/users.sql",
		"internal/pkg/script/migrations/projects.sql",
		"internal/pkg/script/migrations/tasks.sql",
	}

	for _, file := range migrationFiles {
		log.Printf("Running migration file: %s", file)
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Error reading migration file %s: %v", file, err)
			continue
		}

		statements := strings.Split(string(content), ";")
		for i, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			log.Printf("Executing statement %d from %s", i+1, file)
			_, err = db.Exec(stmt)
			if err != nil {
				log.Printf("Error executing statement %d from %s: %v", i+1, file, err)
				log.Printf("Statement was: %s", stmt)
				continue
			}
			log.Printf("Successfully executed statement %d from %s", i+1, file)
		}
		log.Printf("Completed migration file: %s", file)
	}

	log.Println("All migrations completed")
	return nil
}
