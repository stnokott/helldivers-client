package main

import (
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

const _cpuProf = "build/cpu.pprof"

var _duration = 30 * time.Minute

// BenchmarkPProf disguises a pprof run as a benchmark to isolate it from the main loop and profit from its Cleanup function.
func BenchmarkPProf(b *testing.B) {
	b.Logf("will run for %s", _duration.String())

	// start CPU profiling
	fCPU, err := os.Create(_cpuProf)
	if err != nil {
		b.Fatal(err)
	}
	// close file on cleanup
	b.Cleanup(func() {
		if errInner := fCPU.Close(); errInner != nil {
			b.Error(errInner)
		}
	})
	if err := pprof.StartCPUProfile(fCPU); err != nil {
		b.Fatal(err)
	}
	// stop profiling on cleanup
	b.Cleanup(pprof.StopCPUProfile)

	stopChan := make(chan struct{})
	// stop worker on cleanup
	go func() {
		<-time.After(_duration)
		// stop worker
		stopChan <- struct{}{}
	}()
	b.ResetTimer()
	run(stopChan)
}
