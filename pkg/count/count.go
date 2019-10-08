package count

import (
	"math/rand"
	"time"

	"github.com/grindlemire/log"
	"github.com/vrecan/life"
)

// Aggregator is an aggregate interface that we will send counts to
// so they can be aggregated in a thread safe way
type Aggregator interface {
	Aggregate(i int) error
}

// Counter is a counter generating counts (simulating event based output)
type Counter struct {
	id int
	*life.Life
	agg Aggregator
}

// Opt is a configuration function for the functional constructor pattern
type Opt func(*Counter) error

// Agg configures a counter to use a specific aggregator. Handles errors for configuration
func Agg(agg Aggregator) Opt {
	return func(c *Counter) error {
		c.agg = agg
		return nil
	}
}

// NewCounter creates a new counter that generates random counts at random intervals
func NewCounter(id int, opts ...Opt) (c *Counter, err error) {
	// you can specify defaults here that the functional constructor will overwrite
	c = &Counter{
		Life: life.NewLife(),
		id:   id,
	}

	// apply the optional functions and check for errors
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return c, err
		}
	}

	// set the function that will be the "main" of this thread
	c.SetRun(c.run)
	return c, nil
}

// run is the main of the counter goroutine
func (c Counter) run() {
	log.Infof("starting counter %d", c.id)

	// setup our random intervals for generating the count
	randInterval := time.Duration(rand.Intn(10)+1) * time.Second
	generator := time.NewTimer(randInterval)
	for {
		select {
		case <-c.Done:
			// simulate a thread that takes a bit to shut down properly
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
			log.Infof("counter %d successfully shut down", c.id)
			return
		case <-generator.C:
			// generate a random count to apply to our aggregate
			randCount := rand.Intn(10) - 5
			log.Infof("Counter [%d] generating [%d]", c.id, randCount)

			// Send our count to the aggregate
			err := c.agg.Aggregate(randCount)
			if err != nil {
				log.Warnf("Unable to apply aggregate: %v", err)
			}

			randInterval = time.Duration(rand.Intn(10)+1) * time.Second
			generator.Reset(randInterval)
		}
	}
}

// Close satisfies the io.Closer interface and shuts down the thread gracefully
func (c Counter) Close() error {
	log.Infof("shutting down counter %d", c.id)
	// Do other things here if you need to cleanup other stuff
	return c.Life.Close()
}
