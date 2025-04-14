package testhelpers

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"time"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	files, err := migrationFiles()
	pgContainer, err := postgres.Run(ctx, "postgres:17-bullseye",
		postgres.WithInitScripts(files...),
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connectionString,
	}, nil
}

func TruncateAllTables(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || '';
			END LOOP;
		END $$;
	`)

	if err != nil {
		return err
	}

	return nil
}

func migrationFiles() ([]string, error) {
	basePath := filepath.Join("..", "migration")
	files, err := filepath.Glob(filepath.Join(basePath, "*.sql"))

	if err != nil {
		return nil, err
	}

	return files, nil
}
