//this is the page where the tracker page gets implemented

package health

import (
	"sort"
	"sync"
	"time"
)

type RollingTracker struct {
	samples map[string][]time.Duration	//a map for backend, holding a slice of backend samples. every backend gets its own list
	mu sync.Mutex	//to prevent deadlock
	maxSamples int	//bounds
	errormap map[string][]bool
}

func NewRollingTracker(maxSamples int)	*RollingTracker{	//constructor for rolling tracker, returns a new rolling tracker object
	return &RollingTracker{
		samples: make(map[string][]time.Duration),	//initialize map
		maxSamples: maxSamples,
		errormap: make(map[string][]bool),
	}
}

func (t *RollingTracker) Record(backend string, latency time.Duration, isError bool){	//appends to the right backend's slice, drop if maxSamples  exceed
	t.mu.Lock()
	defer t.mu.Unlock()
	t.samples[backend] = append(t.samples[backend], latency)	//Append the latency to t.samples
	t.errormap[backend] = append(t.errormap[backend], isError)	//Add isError to errormap
	if len(t.samples[backend]) > t.maxSamples{		//if slice has grown beyond maxSamples, drop 0the oldest element
		t.samples[backend] = t.samples[backend][1:]		//discards the first element
		t.errormap[backend] = t.errormap[backend][1:]
	}
}

func (t *RollingTracker) Snapshot(backend string)	Snapshot {	//sort the samples list and grabs the value at p95th percentile
	t.mu.Lock()
	defer t.mu.Unlock()
	samples := t.samples[backend]
	if (len(samples) == 0){
		return Snapshot{}	//return nil (a zero-value Sanpshot as the func expects snapshot type return) All three arguments zero.
	}
	sorted := make([]time.Duration, len(samples))	//make a copy of the samples slice, sort them and index into them looking for p95 position
	copy(sorted, samples)
	sort.Slice(sorted, func(i, j int) bool {	//sorted[i] should come before sorted[j] if it's smaller." That's ascending order.
		return sorted[i] < sorted[j]
	})
	idx := int(float64(len(sorted))*0.95)
	errs := t.errormap[backend]
	errorCount := 0
	for _, e := range errs {
		if e {
			errorCount++
		}
	}
	errorRate := float64(errorCount) / float64(len(errs))

	if idx >= len(sorted) {
		idx = len(sorted) - 1	//if we move out of bounds, come back to the last valid index
	}

	return Snapshot{
		P95Latency: sorted[idx],
		SampleCount: len(sorted),
		ErrorRate: errorRate,
	}
}

/*sorted is a slice of latency durations, like [10ms, 250ms, 3ms, 7ms, 18ms, ...] — 
	every recent response's latency, in insertion order (i.e. random).
	After sort.Slice runs, sorted is reordered ascending: [3ms, 7ms, 10ms, 18ms, 250ms, ...].
	That's the point of sorting — once it's ordered, the 95th percentile is literally just "the value at index len * 0.95." 
	Without sorting, the slice is in chronological order and index 95% means nothing useful.
*/