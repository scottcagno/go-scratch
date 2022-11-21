package polling

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type PollerFunc func(t time.Time) error

type Poller struct {
	tick         *time.Ticker
	interval     time.Duration
	running      bool
	fn           PollerFunc
	fnCount      int
	fnStart      time.Time
	fnSinceStart time.Duration
	done         chan bool
}

func NewPoller(d time.Duration, fn PollerFunc) *Poller {
	return &Poller{
		tick:     time.NewTicker(d),
		interval: d,
		fn:       fn,
	}
}

func (bt *Poller) init() chan bool {
	done := make(chan bool)
	go func() {
		var err error
		for {
			select {
			case <-done:
				bt.tick.Stop()
				return
			case t := <-bt.tick.C:
				err = bt.fn(t)
				if err != nil {
					bt.Stop()
					log.Printf("Encountered an error running function, stopping: %q\n", err)
					return
				}
				bt.fnCount++
				bt.fnSinceStart = time.Since(bt.fnStart)
			}
		}
	}()
	return done
}

func (bt *Poller) Start() {
	log.Printf("Poller started...\n")
	log.Printf("Running background func every %s\n", bt.interval)
	if bt.running {
		return
	}
	go func() {
		var err error
		for {
			select {
			case t := <-bt.tick.C:
				err = bt.fn(t)
				if err != nil {
					bt.Stop()
					log.Printf("Encountered an error running function, stopping: %q\n", err)
					return
				}
				bt.fnCount++
				bt.fnSinceStart = time.Since(bt.fnStart)
			}
		}
	}()
	bt.running = true
	bt.fnStart = time.Now()
}

func (bt *Poller) Stop() {
	log.Printf("Poller stopped.\n")
	if bt.running {
		bt.running = false
	}
	bt.done <- true
}

func (bt *Poller) Reset(d time.Duration) {
	log.Printf("Poller reset...\n")
	bt.tick.Reset(d)
	bt.fnStart = time.Now()
	if !bt.running {
		bt.running = true
		log.Printf("Was stopped. Now running background func every %s\n", d)
		return
	} else {
		before := bt.interval
		log.Printf("Was running every %s, now running background func every %s\n", before, d)
	}
	bt.interval = d
}

func (bt *Poller) Destroy() {
	bt.tick = nil
	bt.interval = 0
	bt.running = false
	bt.fn = nil
	runtime.GC()
}

func (bt *Poller) Stats() {
	fmt.Printf("fn has run %d times since it was started %s ago\n", bt.fnCount, bt.fnSinceStart)
}
