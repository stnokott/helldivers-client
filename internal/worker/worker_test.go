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
	"github.com/stnokott/helldivers-client/internal/transform"
)

func mustWorker(t *testing.T) *Worker {
	cfg := config.Get()
	api, err := client.New(cfg, log.Default())
	if err != nil {
		panic(err)
	}
	db, err := db.New(cfg, t.Name(), log.Default())
	if err != nil {
		panic(err)
	}
	return New(api, db, log.Default())
}

func TestWorkerQueryData(t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	worker := mustWorker(t)

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

func testWorkerUpsertDoc[T any](transformer DocTransformer[T], t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	worker := mustWorker(t)

	data := worker.queryData(context.Background())
	upsertDoc(context.Background(), worker, data, transformer)
}

func TestWorkerUpsertPlanets(t *testing.T) {
	testWorkerUpsertDoc(transform.Planets{}, t)
}

func TestWorkerUpsertCampaigns(t *testing.T) {
	testWorkerUpsertDoc(transform.Campaigns{}, t)
}

func TestWorkerUpsertWar(t *testing.T) {
	testWorkerUpsertDoc(transform.War{}, t)
}

func TestWorkerUpserDispatches(t *testing.T) {
	testWorkerUpsertDoc(transform.Dispatches{}, t)
}
