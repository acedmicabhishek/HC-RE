package pkg

import (
	"sync/atomic"
	"time"
)

type Request struct {
	ID         uint64
	Timestamp  time.Time
	CPUWork    int
	IOWait     time.Duration
	MySQLWrite bool
}

type RequestGenerator struct {
	rps    int
	count  uint64
	ticker *time.Ticker
}

func NewRequestGenerator(rps int) *RequestGenerator {
	return &RequestGenerator{
		rps:    rps,
		count:  0,
		ticker: time.NewTicker(time.Second / time.Duration(rps)),
	}
}

func (g *RequestGenerator) Next() *Request {
	<-g.ticker.C
	id := atomic.AddUint64(&g.count, 1)
	return &Request{
		ID:        id,
		Timestamp: time.Now(),
	}
}

func (g *RequestGenerator) Stop() {
	g.ticker.Stop()
}
