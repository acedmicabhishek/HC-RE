package pkg

import "context"

type Dispatcher interface {
	Submit(req *Request) error

	Wait()

	Stop()
}

type Model1Dispatcher struct {
	requests chan *Request
	done     chan struct{}
}

func NewModel1Dispatcher(ctx context.Context) *Model1Dispatcher {
	d := &Model1Dispatcher{
		requests: make(chan *Request, 10000),
		done:     make(chan struct{}),
	}
	go d.run(ctx)
	return d
}

func (d *Model1Dispatcher) Submit(req *Request) error {
	d.requests <- req
	return nil
}

func (d *Model1Dispatcher) Wait() {
	close(d.requests)
	<-d.done
}

func (d *Model1Dispatcher) Stop() {
	d.Wait()
}

func (d *Model1Dispatcher) run(ctx context.Context) {
	for req := range d.requests {
		go func(r *Request) {
			ProcessRequest(r)
		}(req)
	}
	d.done <- struct{}{}
}

type Model2Dispatcher struct {
	workers  int
	requests chan *Request
	done     chan struct{}
}

func NewModel2Dispatcher(ctx context.Context, workers int) *Model2Dispatcher {
	d := &Model2Dispatcher{
		workers:  workers,
		requests: make(chan *Request, workers*10),
		done:     make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		go d.worker(ctx)
	}
	go d.waitForClose()
	return d
}

func (d *Model2Dispatcher) Submit(req *Request) error {
	d.requests <- req
	return nil
}

func (d *Model2Dispatcher) Wait() {
	close(d.requests)
	<-d.done
}

func (d *Model2Dispatcher) Stop() {
	d.Wait()
}

func (d *Model2Dispatcher) worker(ctx context.Context) {
	for req := range d.requests {
		ProcessRequest(req)
	}
}

func (d *Model2Dispatcher) waitForClose() {
	d.done <- struct{}{}
}
