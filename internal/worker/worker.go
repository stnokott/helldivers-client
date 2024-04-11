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
func (w *Worker) Run(interval config.Interval) {
	ticker := time.NewTicker(time.Duration(interval))
	defer ticker.Stop()

	w.log.Printf("worker running every %s", interval.String())

	// this construct forces w.do() to run immediately after starting the ticker.
	// (by default, NewTicker sends the first tick after interval has expired for the 1st time)
	for ; true; <-ticker.C {
		w.do()
	}
}

type docTransformer interface {
	Transform(data transform.APIData) (*db.DocsProvider, error)
}

func (w *Worker) do() {
	w.log.Println("synchronizing")

	var err error
	defer func() {
		if err != nil {
			w.log.Printf("error: %v", err)
		}
		w.log.Println("synchronized")
	}()

	var data transform.APIData
	data, err = w.queryData()
	if err != nil {
		return
	}

	warTransformer := transform.War{}
	if err = w.upsertDoc(data, warTransformer); err != nil {
		return
	}
	planetsTransformer := transform.Planets{}
	if err = w.upsertDoc(data, planetsTransformer); err != nil {
		return
	}
}

func (w *Worker) queryData() (data transform.APIData, err error) {
	data.Planets, err = apiWithTimeout(w.api.Planets, 5*time.Second)
	if err != nil {
		return
	}
	data.WarID, err = apiWithTimeout(w.api.WarID, 1*time.Second)
	if err != nil {
		return
	}
	data.War, err = apiWithTimeout(w.api.War, 5*time.Second)
	if err != nil {
		return
	}
	return
}

func (w *Worker) upsertDoc(data transform.APIData, t docTransformer) error {
	provider, err := t.Transform(data)
	if err != nil {
		return err
	}
	w.db.UpsertDocs(provider, context.TODO())
	return nil
}

func apiWithTimeout[T any](apiFunc func(context.Context) (T, error), timeout time.Duration) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return apiFunc(ctx)
}
