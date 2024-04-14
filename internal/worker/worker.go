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
	Transform(data transform.APIData, errFunc func(error)) *db.DocsProvider[T]
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

	data := w.queryData(ctx)

	if err = w.upsertData(data, ctx); err != nil {
		return
	}
}

func (w *Worker) queryData(ctx context.Context) (data transform.APIData) {
	var err error
	data.WarID, err = w.api.WarID(ctx)
	if err != nil {
		w.log.Printf("failed to query current war ID: %v", err)
	}
	data.War, err = w.api.War(ctx)
	if err != nil {
		w.log.Printf("failed to query current war: %v", err)
	}
	data.Planets, err = w.api.Planets(ctx)
	if err != nil {
		w.log.Printf("failed to query planets: %v", err)
	}
	data.Campaigns, err = w.api.Campaigns(ctx)
	if err != nil {
		w.log.Printf("failed to query campaigns: %v", err)
	}
	data.Dispatches, err = w.api.Dispatches(ctx)
	if err != nil {
		w.log.Printf("failed to query dispatches: %v", err)
	}
	data.Assignments, err = w.api.Assignments(ctx)
	if err != nil {
		w.log.Printf("failed to query assignments: %v", err)
	}
	return
}

func (w *Worker) upsertData(data transform.APIData, ctx context.Context) (err error) {
	warTransformer := transform.War{}
	upsertDoc(w, data, warTransformer, ctx)

	planetsTransformer := transform.Planets{}
	upsertDoc(w, data, planetsTransformer, ctx)

	campaignsTransformer := transform.Campaigns{}
	upsertDoc(w, data, campaignsTransformer, ctx)

	dispatchesTransformer := transform.Dispatches{}
	upsertDoc(w, data, dispatchesTransformer, ctx)

	eventsTransformer := transform.Events{}
	upsertDoc(w, data, eventsTransformer, ctx)

	assignmentsTransformer := transform.Assignments{}
	upsertDoc(w, data, assignmentsTransformer, ctx)
	return
}

func upsertDoc[T any](w *Worker, data transform.APIData, t docTransformer[T], ctx context.Context) {
	provider := t.Transform(data, func(err error) {
		w.log.Printf("error during %T transformation: %v", t, err)
	})
	db.UpsertDocs(w.db, provider, ctx)
}
