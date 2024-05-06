package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

// startProfiling will start CPU profiling. After `duration`, it will stop profiling and will send a signal to `stopChan`.
func startProfiling(duration time.Duration, out string, stopChan chan<- struct{}) error {
	// create output file
	fCPU, err := os.Create(out)
	if err != nil {
		return err
	}
	// start profiling
	if err := pprof.StartCPUProfile(fCPU); err != nil {
		return err
	}

	// wait for duration
	go func() {
		<-time.After(duration)
		fmt.Println("profiling duration reached, stopping worker")
		// stop profiling
		pprof.StopCPUProfile()
		if err := fCPU.Close(); err != nil {
			log.Println(err)
		}
		// stop worker
		stopChan <- struct{}{}
	}()
	return nil
}
