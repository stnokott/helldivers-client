// Package main provides the very simplest main function
package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	workerStopChan := make(chan struct{})

	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan, os.Interrupt)
	go func() {
		s := <-osSignalChan
		fmt.Printf("main loop received %s signal, sending stop signal to worker\n", s.String())
		workerStopChan <- struct{}{}
	}()

	// we separate the run() and main() function so that we can include additional
	// pre- and post-mainloop logic like profiling.
	run(workerStopChan)
}
