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

	// setting timeout to configured interval minus a few seconds to
	// account for rate limit and other overhead
	workTimeout := time.Duration(interval) - 5*time.Second

	// this construct forces w.do() to run immediately after starting the ticker.
	// (by default, NewTicker sends the first tick after interval has expired for the 1st time)
	for ; true; <-ticker.C {
		w.do(workTimeout)
	}
}

type docTransformer[T any] interface {
	Transform(data transform.APIData) (*db.DocsProvider[T], error)
}

func (w *Worker) do(timeout time.Duration) {
	w.log.Println("synchronizing")

	var err error
	defer func() {
		if err != nil {
			w.log.Printf("error: %v", err)
		}
		w.log.Println("synchronized")
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var data transform.APIData
	data, err = w.queryData(ctx)
	if err != nil {
		return
	}
	if err = w.upsertData(data, ctx); err != nil {
		return
	}
}

func (w *Worker) queryData(ctx context.Context) (data transform.APIData, err error) {
	data.Planets, err = w.api.Planets(ctx)
	if err != nil {
		return
	}
	data.WarID, err = w.api.WarID(ctx)
	if err != nil {
		return
	}
	data.War, err = w.api.War(ctx)
	if err != nil {
		return
	}
	data.Campaigns, err = w.api.Campaigns(ctx)
	if err != nil {
		return
	}
	data.Dispatches, err = w.api.Dispatches(ctx)
	if err != nil {
		return
	}
	data.Assignments, err = w.api.Assignments(ctx)
	if err != nil {
		return
	}
	return
}

func (w *Worker) upsertData(data transform.APIData, ctx context.Context) (err error) {
	warTransformer := transform.War{}
	if err = upsertDoc(w, data, warTransformer, ctx); err != nil {
		return
	}
	planetsTransformer := transform.Planets{}
	if err = upsertDoc(w, data, planetsTransformer, ctx); err != nil {
		return
	}
	campaignsTransformer := transform.Campaigns{}
	if err = upsertDoc(w, data, campaignsTransformer, ctx); err != nil {
		return
	}
	dispatchesTransformer := transform.Dispatches{}
	if err = upsertDoc(w, data, dispatchesTransformer, ctx); err != nil {
		return
	}
	eventsTransformer := transform.Events{}
	if err = upsertDoc(w, data, eventsTransformer, ctx); err != nil {
		return
	}
	assignmentsTransformer := transform.Assignments{}
	if err = upsertDoc(w, data, assignmentsTransformer, ctx); err != nil {
		return
	}
	return
}

func upsertDoc[T any](w *Worker, data transform.APIData, t docTransformer[T], ctx context.Context) error {
	provider, err := t.Transform(data)
	if err != nil {
		return err
	}
	db.UpsertDocs(w.db, provider, ctx)
	return nil
}
