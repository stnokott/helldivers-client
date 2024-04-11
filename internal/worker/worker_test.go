// Package worker synchronizes data between API and DB.
//
// It queries the API at a specified interval and merges the results into the DB.
package worker

import (
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

	got, err := worker.queryData()
	if err != nil {
		t.Errorf("Worker.queryData() error = %v, want nil", err)
		return
	}
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

func testWorkerUpsertDoc(transformer docTransformer, t *testing.T) {
	t.Skip("currently skipped until we can disable rate-limiting in API") // TODO
	worker := mustWorker(t)

	data, err := worker.queryData()
	if err != nil {
		t.Skipf("API data not available: %v", err)
		return
	}
	if err := worker.upsertDoc(data, transformer); err != nil {
		t.Errorf("worker.upsertDoc() err = %v, want nil", err)
		return
	}
}

func TestWorkerUpsertPlanets(t *testing.T) {
	testWorkerUpsertDoc(transform.Planets{}, t)
}

func TestWorkerUpsertWar(t *testing.T) {
	testWorkerUpsertDoc(transform.War{}, t)
}
