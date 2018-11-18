package util

import (
	"sync"

	"github.com/DemoHn/apm/mod/instance"
)

// MergeHandles handles multiple instance Events and merge them
// to one Event handler for convenience.
func MergeHandles(ins ...<-chan instance.Event) <-chan instance.Event {
	out := make(chan instance.Event)
	var wg sync.WaitGroup

	output := func(evt <-chan instance.Event) {
		for it := range evt {
			out <- it
		}
		wg.Done()
	}

	wg.Add(len(ins))
	for _, in := range ins {
		go output(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
