//go:build !goverter

// Package worker synchronizes data between API and DB.
//
// It queries the API at a specified interval and merges the results into the DB.
package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	health "github.com/stnokott/healthchecks"

	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/transform"
)

// Worker coordinates communication between API and DB.
type Worker struct {
	api         *client.Client
	db          *db.Client
	healthcheck health.Notifier
	log         *log.Logger
}

// New creates a new Worker instance.
func New(api *client.Client, db *db.Client, cfg *config.Config, logger *log.Logger) (*Worker, error) {
	var healthcheck health.Notifier
	if cfg.HealthchecksURL != "" {
		var err error
		healthcheck, err = health.FromURL(cfg.HealthchecksURL)
		if err != nil {
			return nil, fmt.Errorf("preparing healthcheck: %w", err)
		}
	}

	return &Worker{
		api:         api,
		db:          db,
		healthcheck: healthcheck,
		log:         logger,
	}, nil
}

// Run schedules a new sync job at the specified interval. It is blocking.
func (w *Worker) Run(cron string, stop <-chan struct{}) error {
	// create scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("creating scheduler: %w", err)
	}

	printNextRun := func() {
		if next, errNext := scheduler.Jobs()[0].NextRun(); errNext == nil {
			w.log.Printf("next run at %v", next)
		}
	}

	// create scheduled job
	_, err = scheduler.NewJob(
		gocron.CronJob(cron, false),
		gocron.NewTask(w.do),
		gocron.WithSingletonMode(gocron.LimitModeWait),
		gocron.WithEventListeners(gocron.AfterJobRuns(func(_ uuid.UUID, _ string) {
			printNextRun()
		})),
	)
	if err != nil {
		return fmt.Errorf("creating scheduled job: %w", err)
	}

	scheduler.Start()

	w.log.Printf("started scheduler with cron <%s>", cron)
	printNextRun()

	<-stop
	w.log.Println("received stop signal")
	return scheduler.Shutdown()
}

func (w *Worker) do() {
	w.log.Println("synchronizing")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	w.healthNotify(ctx, healthcheckStart)

	var err error
	defer func() {
		if err != nil {
			w.log.Printf("error: %v", err)
			w.healthNotify(ctx, healthcheckFail)
		} else {
			w.healthNotify(ctx, healthcheckSuccess)
		}
		w.log.Println("synchronized")
	}()

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
	if wars, err = transform.Wars(converter, data); err != nil {
		return
	}
	if campaigns, err = transform.Campaigns(converter, data); err != nil {
		return
	}
	if events, err = transform.Events(converter, data); err != nil {
		return
	}
	if planets, err = transform.Planets(converter, data); err != nil {
		return
	}
	if assignments, err = transform.Assignments(converter, data); err != nil {
		return
	}
	if dispatches, err = transform.Dispatches(converter, data); err != nil {
		return
	}
	if snapshots, err = transform.Snapshot(converter, data); err != nil {
		return
	}

	w.log.Println("merging transformed entities into database")
	// order is important here due to FK constraints
	err = w.db.Merge(ctx, wars, campaigns, events, planets, assignments, dispatches, snapshots)
	return
}
