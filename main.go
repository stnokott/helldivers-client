// Package main provides the very simplest main function
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var (
	pprofDuration = flag.Duration("pprof-duration", 0, "how long to profile for (e.g. 30m)")
	pprofOut      = flag.String("pprof-out", "default.pprof", "where to save the profile")
)

func main() {
	flag.Parse()

	workerStopChan := make(chan struct{})

	// stop worker on SIGINT
	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan, os.Interrupt)
	go func() {
		s := <-osSignalChan
		log.Printf("main loop received %s signal, sending stop signal to worker\n", s.String())
		workerStopChan <- struct{}{}
	}()

	// check if profiling is enabled
	if *pprofDuration != 0 {
		log.Printf("profiling enabled, duration=%s, out=%s\n", pprofDuration.String(), *pprofOut)
		if err := startProfiling(*pprofDuration, *pprofOut, workerStopChan); err != nil {
			log.Fatalln(err)
		}
	}

	// we separate the run() and main() function so that we can include additional
	// pre- and post-mainloop logic like profiling.
	run(workerStopChan)
}
