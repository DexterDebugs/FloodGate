package ratelimit

import (
	"github.com/DexterDebugs/FloodGate/internal/control"
	"github.com/DexterDebugs/FloodGate/internal/health"
)

type Adaptive struct {
	baseLimiter Limiter
	tracker health.Tracker
	pid *control.PIDController
	baseLimit int
	targetP95 float64
	backendName string
}

func NewAdaptive(base Limiter, tracker health.Tracker, pid *control.PIDController, baseLimit int, targetP95 float64, backendName string) *Adaptive {
	return &Adaptive{
		baseLimiter: base,
		tracker: tracker,
		pid: pid,
		baseLimit: baseLimit,
		targetP95: targetP95,
		backendName: backendName,
	}	
}

func (a *Adaptive) Allow(clientID, route string) bool {
	snap := a.tracker.Snapshot(a.backendName)		//knowing the current health state of this adaptive's backend. tracker returns a snapshot with 3 fields
		//p95latency, errorrate, samplecount
	if snap.SampleCount == 0 {		//if tracker has zero samples for this backend, the snapshots other field remain zero
		return a.baseLimiter.Allow(clientID, route)
	}

	currentP95Ms := float64(snap.P95Latency.Milliseconds())	//converting from nano to milli(int64) and then int64 to float64
	currentError := currentP95Ms - a.targetP95		//the error signal used to pipe into the input to PID, if +ve - tighten, else loosen



	output := a.pid.Compute(currentError)	//Hand the error to the controller and get back the adjustment signal, the PID output

	newLimit := a.baseLimit - int(output)		//acts like floor function, guards negative
	if newLimit < 1 { newLimit = 1 }
	if newLimit > 2 * a.baseLimit { newLimit =  2* a.baseLimit}			//acts like ceiling function 
	a.baseLimiter.SetLimit(newLimit)		//tell the wrapped limiter(fw or sw) its new ceiling


	return a.baseLimiter.Allow(clientID, route)		
}

func (a *Adaptive) SetLimit(limit int) {
    a.baseLimit = limit
}