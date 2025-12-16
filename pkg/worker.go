package pkg

import (
	"time"
)

func ProcessRequest(req *Request) {

	if req.CPUWork > 0 {
		CPUBound(req.CPUWork)
	}

	if req.IOWait > 0 {
		time.Sleep(req.IOWait)
	}

	if req.MySQLWrite {

	}
}

func CPUBound(iterations int) {
	sum := 0
	for i := 0; i < iterations; i++ {
		sum += i
	}
	_ = sum
}
