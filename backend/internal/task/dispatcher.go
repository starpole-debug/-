package task

import (
	"context"
)

// Dispatcher provides a tiny worker pool for fire-and-forget events like notification fan out.
type Dispatcher struct {
	tasks chan func(context.Context)
	quit  chan struct{}
}

// NewDispatcher spins up a worker that processes background closures sequentially.
func NewDispatcher(buffer int) *Dispatcher {
	if buffer <= 0 {
		buffer = 8
	}
	d := &Dispatcher{
		tasks: make(chan func(context.Context), buffer),
		quit:  make(chan struct{}),
	}
	go d.loop()
	return d
}

// Dispatch pushes a new asynchronous job.
func (d *Dispatcher) Dispatch(fn func(context.Context)) {
	select {
	case d.tasks <- fn:
	default:
		// Drop tasks when queue full to keep service responsive; acceptable for MVP notifications.
	}
}

// Stop gracefully stops the worker.
func (d *Dispatcher) Stop() {
	close(d.quit)
}

func (d *Dispatcher) loop() {
	ctx := context.Background()
	for {
		select {
		case task := <-d.tasks:
			if task != nil {
				task(ctx)
			}
		case <-d.quit:
			return
		}
	}
}
