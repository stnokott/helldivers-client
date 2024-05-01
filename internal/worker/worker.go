//go:build !goverter

// TODO: remove

// Package worker synchronizes data between API and DB.
//
// It queries the API at a specified interval and merges the results into the DB.
package worker

import (
	"context"
	"log"
	"time"

	"github.com/stnokott/helldivers-client/internal/client"
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
func (w *Worker) Run(interval time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	w.log.Printf("worker running every %s", interval.String())

	// setting timeout to configured interval minus a few seconds to
	// account for rate limit and other overhead
	workTimeout := time.Duration(interval) - 5*time.Second

	// this construct forces w.do() to run immediately after starting the ticker.
	// (by default, NewTicker sends the first tick after interval has expired for the 1st time)
	for {
		w.do(workTimeout)
		select {
		case <-ticker.C:
			continue
		case <-stop:
			w.log.Println("received stop signal")
			return
		}
	}
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

	if err = w.mergeData(ctx, data); err != nil {
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
	data.Campaigns, err = w.api.Campaigns(ctx)
	if err != nil {
		w.log.Printf("failed to query campaigns: %v", err)
	}
	data.Planets, err = w.api.Planets(ctx)
	if err != nil {
		w.log.Printf("failed to query planets: %v", err)
	}
	data.Assignments, err = w.api.Assignments(ctx)
	if err != nil {
		w.log.Printf("failed to query assignments: %v", err)
	}
	data.Dispatches, err = w.api.Dispatches(ctx)
	if err != nil {
		w.log.Printf("failed to query dispatches: %v", err)
	}
	return
}

func (w *Worker) mergeData(ctx context.Context, data transform.APIData) (err error) {
	w.log.Println("transforming API responses")
	var wars, events, planets, campaigns, assignments, dispatches, snapshots []db.EntityMerger

	converter := &transform.ConverterImpl{}
	if wars, err = transform.Wars(data); err != nil {
		return
	}
	if campaigns, err = transform.Campaigns(data); err != nil {
		return
	}
	if events, err = transform.Events(data); err != nil {
		return
	}
	if planets, err = transform.Planets(data); err != nil {
		return
	}
	if assignments, err = transform.Assignments(converter, data); err != nil {
		return
	}
	if dispatches, err = transform.Dispatches(converter, data); err != nil {
		return
	}
	if snapshots, err = transform.Snapshot(data); err != nil {
		return
	}

	w.log.Println("merging transformed entities into database")
	// order is important here due to FK constraints
	err = w.db.Merge(ctx, wars, campaigns, events, planets, assignments, dispatches, snapshots)
	return
}
