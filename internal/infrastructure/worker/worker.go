package worker

import (
	"context"
	"database/sql"

	"github.com/ahhhmadtlz/series_reader_backend/internal/observability/logger"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
)

func NewRiverClient(db *sql.DB, imageWorker *ImageProcessingWorker) (*river.Client[*sql.Tx], error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, imageWorker)

	riverClient, err := river.NewClient(riverdatabasesql.New(db), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 5},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, err
	}

	return riverClient, nil
}

func StartWorker(ctx context.Context, client *river.Client[*sql.Tx]) error {
	if err := client.Start(ctx); err != nil {
		return err
	}
	logger.Info("River worker started")
	return nil
}

func StopWorker(ctx context.Context, client *river.Client[*sql.Tx]) {
	if err := client.Stop(ctx); err != nil {
		logger.Error("failed to stop River worker", "error", err)
		return
	}
	logger.Info("River worker stopped")
}