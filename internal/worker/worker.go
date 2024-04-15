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

// DocTransformer provides a means of converting API data into structs ready for passing to the MongoDB driver.
type DocTransformer[T any] interface {
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

	if err = w.upsertData(ctx, data); err != nil {
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

func (w *Worker) upsertData(ctx context.Context, data transform.APIData) (err error) {
	warTransformer := transform.War{}
	upsertDoc(ctx, w, data, warTransformer)

	planetsTransformer := transform.Planets{}
	upsertDoc(ctx, w, data, planetsTransformer)

	campaignsTransformer := transform.Campaigns{}
	upsertDoc(ctx, w, data, campaignsTransformer)

	dispatchesTransformer := transform.Dispatches{}
	upsertDoc(ctx, w, data, dispatchesTransformer)

	eventsTransformer := transform.Events{}
	upsertDoc(ctx, w, data, eventsTransformer)

	assignmentsTransformer := transform.Assignments{}
	upsertDoc(ctx, w, data, assignmentsTransformer)

	snapshotsTransformer := transform.Snapshots{}
	upsertDoc(ctx, w, data, snapshotsTransformer)
	return
}

func upsertDoc[T any](ctx context.Context, w *Worker, data transform.APIData, t DocTransformer[T]) {
	provider := t.Transform(data, func(err error) {
		w.log.Printf("error during %T transformation: %v", t, err)
	})
	db.UpsertDocs(ctx, w.db, provider)
}
