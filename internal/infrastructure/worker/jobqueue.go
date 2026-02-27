package worker

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/imageprocessing/param"
	"github.com/riverqueue/river"
)

type JobQueue interface {
	Enqueue(ctx context.Context, args param.ProcessImageArgs) error
}

type RiverJobQueue struct {
	client *river.Client[*sql.Tx]
}

func NewRiverJobQueue(client *river.Client[*sql.Tx]) *RiverJobQueue {
	return &RiverJobQueue{client: client}
}

func (q *RiverJobQueue) Enqueue(ctx context.Context, args param.ProcessImageArgs) error {
	_, err := q.client.Insert(ctx, args, nil)
	return err
}

// LazyJobQueue is a JobQueue whose underlying client is set after construction.
// This breaks the circular dependency between setupServices (which builds the
// imageWorker) and main (which needs the imageWorker to build the River client).
//
// Set must be called exactly once before the HTTP server starts accepting
// requests. If Enqueue is called before Set, it returns a clear error instead
// of panicking on a nil pointer.
type LazyJobQueue struct {
	queue unsafe.Pointer // stores *RiverJobQueue via atomic
}

func NewLazyJobQueue() *LazyJobQueue {
	return &LazyJobQueue{}
}

func (l *LazyJobQueue) Set(q *RiverJobQueue) {
	atomic.StorePointer(&l.queue, unsafe.Pointer(q))
}

func (l *LazyJobQueue) Enqueue(ctx context.Context, args param.ProcessImageArgs) error {
	ptr := atomic.LoadPointer(&l.queue)
	if ptr == nil {
		return fmt.Errorf("job queue not initialized: Set must be called before Enqueue")
	}
	return (*RiverJobQueue)(ptr).Enqueue(ctx, args)
}