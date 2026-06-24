//tracking latency samples per backend and computing p95 on demand. this is what it does. This is just a Design interface

package health

import "time"


//What gets returned from a snapshot
type Snapshot struct {
	P95Latency time.Duration
	ErrorRate float64	//fraction of recent requests that failed
	SampleCount int		//so it knows it has enough data to trust 
}

//that snapshot converted into JSON format for readability
type SnapshotJSON struct {	
	P95LatencyMS int `json:"p95_latency_ms"`
	ErrorRate float64 `json:"error_rate"`
	SampleCount int `json:"sample_count"`
}
//the function which converts 
func (s Snapshot) ToJSON() SnapshotJSON {
	return SnapshotJSON{
		P95LatencyMS: int(s.P95Latency.Milliseconds()),
		ErrorRate: s.ErrorRate,
		SampleCount: s.SampleCount,
	}
}

//Tracker contact 
type Tracker interface {
	Record(backend string, latency time.Duration, isError bool)	//returns nothing. This is intentional — recording is fire-and-forget;
	//  if it fails, we don't want the request to fail.
	Snapshot(backend string) Snapshot	//returns a snapshot method with P95 Latency computed
}