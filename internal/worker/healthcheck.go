package worker

import "context"

type healthcheckType int

const (
	healthcheckStart healthcheckType = iota
	healthcheckSuccess
	healthcheckFail
)

func (w *Worker) healthNotify(ctx context.Context, t healthcheckType) {
	if w.healthcheck == nil {
		return
	}
	var (
		healthFunc func(context.Context) error
		healthName string
	)
	switch t {
	case healthcheckStart:
		healthFunc = w.healthcheck.Start
		healthName = "start"
	case healthcheckSuccess:
		healthFunc = w.healthcheck.Success
		healthName = "success"
	case healthcheckFail:
		healthFunc = w.healthcheck.Fail
		healthName = "fail"
	}

	w.log.Println("signalling healthcheck", healthName)
	if err := healthFunc(ctx); err != nil {
		w.log.Println("WARN: failed to signal:", err)
	}
}
