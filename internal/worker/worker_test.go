// Package worker synchronizes data between API and DB.
//
// It queries the API at a specified interval and merges the results into the DB.
package worker

import (
	"context"
	"log"
	"testing"

	"github.com/stnokott/helldivers-client/internal/client"
	"github.com/stnokott/helldivers-client/internal/config"
	"github.com/stnokott/helldivers-client/internal/db"
)

func mustWorker() *Worker {
	cfg := config.MustGet()
	api, err := client.New(cfg, log.Default())
	if err != nil {
		panic(err)
	}
	db, err := db.New(cfg, log.Default())
	if err != nil {
		panic(err)
	}
	return New(api, db, log.Default())
}

func TestWorkerQueryData(t *testing.T) {
	worker := mustWorker()

	got := worker.queryData(context.Background())
	if got.Planets == nil {
		t.Error("Worker.queryData().Planets = nil, want non-nil")
		return
	}
	if got.WarID == nil {
		t.Error("Worker.queryData().WarID = nil, want non-nil")
		return
	}
	if got.War == nil {
		t.Error("Worker.queryData().War = nil, want non-nil")
		return
	}
}
