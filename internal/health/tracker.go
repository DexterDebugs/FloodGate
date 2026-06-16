package health

import "time"


//What gets returned from a snapshot
type Snapshot struct {
	P95Latency time.Duration
	ErrorRate float64	//fraction of recent requests that failed
	SampleCount int		//so it knows it has enough data to trust 
}

//Tracker contact 
type Tracker interface {
	Record(backend string, latency time.Duration, isError bool)	//returns nothing. This is intentional — recording is fire-and-forget;
	//  if it fails, we don't want the request to fail.
	Snapshot(backend string) Snapshot
}