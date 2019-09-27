package aggregate

import (
	"fmt"

	"github.com/grindlemire/log"
	"github.com/vrecan/life"
)

// Sum sums up all the incoming counts
type Sum struct {
	*life.Life
	sum   int
	input chan int
}

// NewSum creates a new sum thread
func NewSum() (s *Sum, err error) {
	s = &Sum{
		Life:  life.NewLife(),
		sum:   0,
		input: make(chan int, 10), // create an input buffer size of 10. This is often configurable
	}
	s.SetRun(s.run)
	return s, nil
}

// run will run the goRoutine. It is the "main" of this thread
// note I am using a pointer receiver because we are mutating s.Sum here
func (s *Sum) run() {
	log.Info("starting sum thread")
	for {
		select {
		case <-s.Done:
			log.Info("successfully shut down sum thread")
			return
		case i := <-s.input:
			log.Infof("aggregate summing [%d]. sum is [%d]->[%d]", i, s.sum, s.sum+i)
			s.sum += i
		}
	}
}

// Apply applies an integer to the sum in a thread safe manner. This is the only function used by other
// threads (other than close) and is the only function in the interface
func (s Sum) Apply(i int) (err error) {
	select {
	case s.input <- i:
		return nil
	default:
		return fmt.Errorf("unable to apply sum. Channel full [%d]", len(s.input))
	}
}

// Close shuts down the sum thread
func (s Sum) Close() error {
	log.Info("shutting down sum aggregator")
	return s.Life.Close()
}
