package alb_log_to_json

import (
	"testing"
)

func TestAggregator_BasicGrouping(t *testing.T) {
	agg := NewAggregator("domain_name")

	agg.Add(&ALBLog{DomainName: "a.com"})
	agg.Add(&ALBLog{DomainName: "b.com"})
	agg.Add(&ALBLog{DomainName: "a.com"})
	agg.Add(&ALBLog{DomainName: "a.com"})
	agg.Add(&ALBLog{DomainName: "b.com"})

	results := agg.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	// Sorted by count descending
	if results[0].Key != "a.com" || results[0].Count != 3 {
		t.Errorf("expected first group {a.com, 3}, got {%s, %d}", results[0].Key, results[0].Count)
	}
	if results[1].Key != "b.com" || results[1].Count != 2 {
		t.Errorf("expected second group {b.com, 2}, got {%s, %d}", results[1].Key, results[1].Count)
	}
}

func TestAggregator_NilPointerField(t *testing.T) {
	agg := NewAggregator("target_status_code")

	code200 := 200
	agg.Add(&ALBLog{TargetStatusCode: &code200})
	agg.Add(&ALBLog{TargetStatusCode: nil})
	agg.Add(&ALBLog{TargetStatusCode: nil})

	results := agg.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	// nil group has count 2, should be first
	if results[0].Key != "<nil>" || results[0].Count != 2 {
		t.Errorf("expected first group {<nil>, 2}, got {%s, %d}", results[0].Key, results[0].Count)
	}
}

func TestAggregator_SingleGroup(t *testing.T) {
	agg := NewAggregator("request_method")

	agg.Add(&ALBLog{RequestMethod: "GET"})
	agg.Add(&ALBLog{RequestMethod: "GET"})

	results := agg.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 group, got %d", len(results))
	}
	if results[0].Key != "GET" || results[0].Count != 2 {
		t.Errorf("expected {GET, 2}, got {%s, %d}", results[0].Key, results[0].Count)
	}
}

func TestAggregator_Empty(t *testing.T) {
	agg := NewAggregator("domain_name")
	results := agg.Results()
	if len(results) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(results))
	}
}
