package _main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"
)

// v1: https://go.dev/play/p/nKtSKjKxK_o
// v2: https://go.dev/play/p/1ottKCJOUKQ

var myBgFunc = func(t time.Time) error {
	fmt.Printf("running my background func: %ds (elapsed)\n", t.Second())
	if t.Second() > 10 {
		return errors.New("something weird happened")
	}
	return nil
}

func init() {
	runtime.GOMAXPROCS(1)
}

func main() {

	b := NewBgTicker(2*time.Second, myBgFunc)
	b.Start()
	time.Sleep(8 * time.Second)
	b.Stats()
	b.Reset(1 * time.Second)
	time.Sleep(4 * time.Second)
	b.Stop()
	time.Sleep(2 * time.Second)
	b.Reset(3 * time.Second)
	time.Sleep(15 * time.Second)
	b.Stop()
}

type BgFn func(t time.Time) error

type BgTicker struct {
	tick         *time.Ticker
	interval     time.Duration
	running      bool
	fn           BgFn
	fnCount      int
	fnStart      time.Time
	fnSinceStart time.Duration
}

func NewBgTicker(d time.Duration, fn BgFn) *BgTicker {
	return &BgTicker{
		tick:     time.NewTicker(d),
		interval: d,
		fn:       fn,
	}
}

func (bt *BgTicker) Start() {
	log.Printf("BgTicker started...\n")
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

func (bt *BgTicker) Stop() {
	log.Printf("BgTicker stopped.\n")
	if bt.running {
		bt.running = false
	}
	bt.tick.Stop()
}

func (bt *BgTicker) Reset(d time.Duration) {
	log.Printf("BgTicker reset...\n")
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

func (bt *BgTicker) Destroy() {
	bt.tick = nil
	bt.interval = 0
	bt.running = false
	bt.fn = nil
	runtime.GC()
}

func (bt *BgTicker) Stats() {
	fmt.Printf("fn has run %d times since it was started %s ago\n", bt.fnCount, bt.fnSinceStart)
}
