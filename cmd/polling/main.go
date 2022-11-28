package main

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func main() {

	p := NewPoller(500*time.Millisecond, nil)
	// fmt.Println(p)
	fmt.Println(">> sleeping for 1 second...")
	time.Sleep(1 * time.Second)

	fmt.Println(">> starting")
	p.Start()
	// fmt.Println(p)
	fmt.Println(">> sleeping for 2 seconds... (should see 4 events triggered)")
	time.Sleep(2 * time.Second)
	// should see 4 events triggered

	p.Stop()

	fmt.Println(">> resetting to 2 second intervals")
	p.Reset(2 * time.Second)
	fmt.Println(">> sleeping for 4 seconds... (should see 2 events triggered)")
	time.Sleep(4 * time.Second)
	// should see 2 events triggered

	p.Stop()
	// fmt.Println(p)

}

type PollerFunc func(t time.Time)

func defaultPollerFunc(t time.Time) {
	fmt.Println("Tick at", t)
}

type Poller struct {
	ticker  *time.Ticker
	fn      PollerFunc
	done    chan bool
	running bool
}

func NewPoller(d time.Duration, fn PollerFunc) *Poller {
	if fn == nil {
		fn = defaultPollerFunc
	}
	p := &Poller{
		ticker: time.NewTicker(d),
		fn:     fn,
		done:   make(chan bool),
	}
	return p
}

func (p *Poller) Start() {
	p.startEventLoop()
	log.Printf("Poller started\n")
}

func (p *Poller) startEventLoop() {
	if p.running {
		return
	}
	p.running = true
	go func() {
		for {
			select {
			case <-p.done:
				return
			case t := <-p.ticker.C:
				p.fn(t)
			}
		}
	}()
}

func (p *Poller) Reset(d time.Duration) {
	p.resetEventLoop(d)
	log.Printf("Poller reset to %s intervals\n", d)
}

func (p *Poller) resetEventLoop(d time.Duration) {
	// stop the event loop
	p.stopEventLoop()
	// reset the ticker with a new interval
	p.ticker.Reset(d)
	// start the event loop again
	p.startEventLoop()
}

func (p *Poller) Stop() {
	p.ticker.Stop()
	p.stopEventLoop()
	log.Printf("Poller stopped.\n")
}

func (p *Poller) stopEventLoop() {
	if p.running {
		// kill the existing event loop
		p.done <- true
		// make sure running is false
		p.running = false
	}
}

func (p *Poller) String() string {
	return fmt.Sprintf("fn=%q, running=%v\n", GetFuncName(p.fn), p.running)
}

func GetFuncName(fn any) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if n := strings.IndexByte(name, '.'); n > 0 {
		name = name[n+1:]
	}
	return name
}
