package polling

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var myPollerFunc = func(t time.Time) error {
	fmt.Printf("running my background func: %ds (elapsed)\n", t.Second())
	if t.Second() > 10 {
		return errors.New("something weird happened")
	}
	return nil
}

func TestPoller_All(t *testing.T) {

	fmt.Printf(">>> Create new 2s poller.\n")
	b := NewPoller(2*time.Second, myPollerFunc)

	fmt.Printf(">>> Sleeping for 3 seconds...\n")
	time.Sleep(3 * time.Second)

	fmt.Printf(">>> Starting poller.\n")
	b.Start()

	fmt.Printf(">>> Sleeping for 8 seconds...\n")
	time.Sleep(8 * time.Second)

	fmt.Printf(">>> Getting poller stats.\n")
	b.Stats()

	fmt.Printf(">>> Resetting poller.\n")
	b.Reset(1 * time.Second)

	fmt.Printf(">>> Sleeping for 4 seconds...\n")
	time.Sleep(4 * time.Second)

	b.Stop()

	fmt.Printf(">>> Sleeping for 2 seconds...\n")
	time.Sleep(2 * time.Second)

	b.Reset(3 * time.Second)

	fmt.Printf(">>> Sleeping for 15 seconds...\n")
	time.Sleep(15 * time.Second)

	fmt.Printf(">>> Stopping poller.\n")
	b.Stop()
}
