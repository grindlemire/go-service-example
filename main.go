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

	// create 10 pools of go routines to show signal propogation
	for i := 0; i < 10; i++ {
		count, err := count.NewCounter(i, s)
		if err != nil {
			log.Fatal(err)
		}
		count.Start()
		goRoutines = append(goRoutines, count)
	}

	d.WaitForDeath(goRoutines...)
	log.Info("successfully shutdown all go routines")
}