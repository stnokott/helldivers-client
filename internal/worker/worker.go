// Package worker synchronizes data between API and DB.
//
// It queries the API at a specified interval and merges the results into the DB.
package worker

import (
	"context"
	"log"
	"time"

	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/transform"
)

// Worker coordinates communication between API and DB.
type Worker struct {
	api *client.Client
	db  *db.Client
	log *log.Logger
}

// New creates a new Worker instance.
func New(api *client.Client, db *db.Client, logger *log.Logger) *Worker {
	return &Worker{
		api: api,
		db:  db,
		log: logger,
	}
}

// Run schedules a new sync job at the specified interval. It is blocking.
func (w *Worker) Run(interval config.WorkerInterval) {
	ticker := time.NewTicker(time.Duration(interval))
	defer ticker.Stop()

	w.log.Printf("worker running every %s", interval.String())

	// this construct forces w.do() to run immediately after starting the ticker.
	// (by default, NewTicker sends the first tick after interval has expired for the 1st time)
	for ; true; <-ticker.C {
		w.do()
	}
}

type docTransformer[T any] interface {
	Request(api *client.Client, ctx context.Context) (T, error)
	Transform(data T) (db.DocProvider, error)
}

func processDoc[T any](w *Worker, t docTransformer[T]) error {
	data, err := t.Request(w.api, context.TODO())
	if err != nil {
		return err
	}
	provider, err := t.Transform(data)
	if err != nil {
		return err
	}
	inserted, err := w.db.UpsertDoc(provider, context.TODO())
	if err != nil {
		return err
	}
	if inserted {
		w.log.Printf("new %s document added. ID=%d", provider.CollectionName(), provider.DocID())
	}
	return nil
}

func (w *Worker) do() {
	w.log.Println("synchronizing")

	var err error
	defer func() {
		if err != nil {
			w.log.Printf("error: %v", err)
		} else {
			w.log.Println("synchronized")
		}
	}()

	warTransformer := transform.War{}
	if err = processDoc(w, warTransformer); err != nil {
		return
	}
}
