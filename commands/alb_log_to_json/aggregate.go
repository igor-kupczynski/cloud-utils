package alb_log_to_json

import (
	"fmt"
	"sort"
)

// GroupResult represents a single group in --group-by output.
type GroupResult struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// Aggregator counts records grouped by a field value.
type Aggregator struct {
	extract func(*ALBLog) interface{}
	counts  map[string]int
}

// NewAggregator creates an Aggregator for the given FieldRegistry field.
func NewAggregator(field string) *Aggregator {
	return &Aggregator{
		extract: FieldRegistry[field],
		counts:  make(map[string]int),
	}
}

// Add extracts the field value from the log entry and increments its counter.
func (a *Aggregator) Add(log *ALBLog) {
	key := fmt.Sprintf("%v", a.extract(log))
	a.counts[key]++
}

// Results returns the aggregated counts sorted by count descending.
func (a *Aggregator) Results() []GroupResult {
	results := make([]GroupResult, 0, len(a.counts))
	for key, count := range a.counts {
		results = append(results, GroupResult{Key: key, Count: count})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})
	return results
}
