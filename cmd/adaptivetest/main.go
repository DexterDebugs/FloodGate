package main

import (
	"fmt"
	"time"

	"github.com/DexterDebugs/FloodGate/internal/control"
	"github.com/DexterDebugs/FloodGate/internal/health"
	"github.com/DexterDebugs/FloodGate/internal/ratelimit"
)

type FakeTracker struct {
	snap health.Snapshot
}

type FakeBaseLimiter struct {
	capturedLimit int
}

func (f *FakeTracker)	Record(backend string, latency time.Duration, isError bool){

}

func (f *FakeTracker)	Snapshot(backend string)	health.Snapshot {
	return f.snap
}

func (f *FakeBaseLimiter)	SetLimit(limit int) {
	f.capturedLimit = limit
}

func (f *FakeBaseLimiter)	Allow(clientID, route string)	bool {
	return true
}

func main() {
	tracker := &FakeTracker{
		snap: health.Snapshot{
			P95Latency: 100 * time.Millisecond,
			SampleCount: 50,
		},
	}
	base := &FakeBaseLimiter{}

	pid := control.NewPIDController(0.1,0,0)

	adaptive := ratelimit.NewAdaptive(base, tracker, pid, 10, 100.0, "test")

	p95Values := []int{100, 200, 150, 50}
	for i, p95 := range p95Values {
		tracker.snap.P95Latency = time.Duration(p95) * time.Millisecond
		adaptive.Allow("test-client", "/test")
		fmt.Printf("Tick %d: p95=%dms -> limit=%d\n", i+1, p95, base.capturedLimit)
	}

}