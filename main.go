package main

import (
	"io"
	"syscall"

	"github.com/grindlemire/go-service-example/pkg/aggregate"
	"github.com/grindlemire/go-service-example/pkg/count"
	"github.com/grindlemire/log"
	"github.com/vrecan/death"
)

const children = 10

func main() {
	log.Init(log.Default)

	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	goRoutines := []io.Closer{}

	// Start our sum aggregate in one go routine
	s, err := aggregate.NewSum()
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
	goRoutines = append(goRoutines, s)

	// create 10 threads generating random counts at random intervals
	for i := 0; i < 10; i++ {
		count, err := count.NewCounter(i, s)
		if err != nil {
			log.Fatal(err)
		}
		count.Start()
		goRoutines = append(goRoutines, count)
	}

	err = d.WaitForDeath(goRoutines...)
	if err != nil {
		log.Fatalf("failed to cleanly shut down all go routines: %v", err)
	}

	log.Info("successfully shutdown all go routines")
}
