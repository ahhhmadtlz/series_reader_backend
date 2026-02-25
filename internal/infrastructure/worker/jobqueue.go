package worker

import (
	"context"
	"database/sql"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	"github.com/riverqueue/river"
)

type JobQueue interface {
	Enqueue(ctx context.Context, args param.ProcessImageArgs)error
}


type RiverJobQueue struct {
	client *river.Client[*sql.Tx]
}

func (q *RiverJobQueue) SetClient(client *river.Client[*sql.Tx]) {
    q.client = client
}


func NewRiverJobQueue(client *river.Client[*sql.Tx]) *RiverJobQueue {
	return &RiverJobQueue{client: client}
}


func (q *RiverJobQueue) Enqueue(ctx context.Context,args param.ProcessImageArgs)error{
_, err := q.client.Insert(ctx, args, nil)
	return err
}
