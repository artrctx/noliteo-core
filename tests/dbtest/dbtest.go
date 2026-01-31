package dbtest

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	cfg "github.com/artrctx/quoin-core/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartPostgresContainer(config *cfg.DatabaseConfig) (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	config.Database = "testdb"
	config.Username = "testusr"
	config.Password = "password"

	ctx := context.Background()
	dbContainer, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase(config.Database),
		postgres.WithUsername(config.Username),
		postgres.WithPassword(config.Password),
		// mute test container log
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
		testcontainers.WithAdditionalWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(ctx)
	if err != nil {
		return dbContainer.Terminate, err
	}
	config.Host = dbHost

	dbPort, err := dbContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}
	config.Port = dbPort.Port()
	config.SSLMode = "disable"
	config.Schema = "public"

	return dbContainer.Terminate, nil
}

func RunWithPostgres(m *testing.M, cfg *cfg.DatabaseConfig) {
	teardown, err := StartPostgresContainer(cfg)
	if err != nil {
		log.Fatalf("could not start postgres test container: %v", err)
	}

	m.Run()

	if err := teardown(context.Background()); err != nil {
		log.Fatalf("could not teardown postgres test container: %v", err)
	}
}
