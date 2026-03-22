package alb_log_to_json

import (
	"testing"
	"time"
)

func TestStatusMatcher_Exact(t *testing.T) {
	m, err := parseStatusPattern("502")
	if err != nil {
		t.Fatal(err)
	}

	code502 := 502
	code200 := 200
	if !m.matches(&code502) {
		t.Error("expected 502 to match")
	}
	if m.matches(&code200) {
		t.Error("expected 200 not to match")
	}
	if m.matches(nil) {
		t.Error("expected nil not to match")
	}
}

func TestStatusMatcher_Class(t *testing.T) {
	m, err := parseStatusPattern("5xx")
	if err != nil {
		t.Fatal(err)
	}

	code502 := 502
	code503 := 503
	code200 := 200
	if !m.matches(&code502) {
		t.Error("expected 502 to match 5xx")
	}
	if !m.matches(&code503) {
		t.Error("expected 503 to match 5xx")
	}
	if m.matches(&code200) {
		t.Error("expected 200 not to match 5xx")
	}
	if m.matches(nil) {
		t.Error("expected nil not to match 5xx")
	}
}

func TestStatusMatcher_Invalid(t *testing.T) {
	_, err := parseStatusPattern("abc")
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
	_, err = parseStatusPattern("0xx")
	if err == nil {
		t.Error("expected error for 0xx")
	}
}

func TestFilter_ElbStatus(t *testing.T) {
	m, _ := parseStatusPattern("502")
	f := &Filter{ElbStatus: &m}

	code502 := 502
	code200 := 200
	if !f.Matches(&ALBLog{ElbStatusCode: &code502}) {
		t.Error("expected 502 to match")
	}
	if f.Matches(&ALBLog{ElbStatusCode: &code200}) {
		t.Error("expected 200 not to match")
	}
}

func TestFilter_Domain(t *testing.T) {
	f := &Filter{Domain: "example.com"}

	if !f.Matches(&ALBLog{DomainName: "www.example.com"}) {
		t.Error("expected substring match")
	}
	if f.Matches(&ALBLog{DomainName: "other.com"}) {
		t.Error("expected no match")
	}
}

func TestFilter_Method(t *testing.T) {
	f := &Filter{Method: "GET"}

	if !f.Matches(&ALBLog{RequestMethod: "GET"}) {
		t.Error("expected GET to match")
	}
	if !f.Matches(&ALBLog{RequestMethod: "get"}) {
		t.Error("expected case-insensitive match")
	}
	if f.Matches(&ALBLog{RequestMethod: "POST"}) {
		t.Error("expected POST not to match")
	}
}

func TestFilter_Path(t *testing.T) {
	f := &Filter{Path: "/api/"}

	if !f.Matches(&ALBLog{RequestURL: "http://example.com/api/health"}) {
		t.Error("expected path substring match")
	}
	if f.Matches(&ALBLog{RequestURL: "http://example.com/other"}) {
		t.Error("expected no match")
	}
}

func TestFilter_TimeRange(t *testing.T) {
	after, _ := time.Parse(time.RFC3339, "2024-01-15T12:00:00Z")
	before, _ := time.Parse(time.RFC3339, "2024-01-15T12:01:00Z")
	f := &Filter{After: &after, Before: &before}

	if !f.Matches(&ALBLog{Time: "2024-01-15T12:00:30.000000Z"}) {
		t.Error("expected time within range to match")
	}
	if f.Matches(&ALBLog{Time: "2024-01-15T11:59:59.000000Z"}) {
		t.Error("expected time before range not to match")
	}
	if f.Matches(&ALBLog{Time: "2024-01-15T12:01:01.000000Z"}) {
		t.Error("expected time after range not to match")
	}
	// Boundary: at exact after time should match (>=)
	if !f.Matches(&ALBLog{Time: "2024-01-15T12:00:00.000000Z"}) {
		t.Error("expected time at after boundary to match")
	}
	// Boundary: at exact before time should match (<=)
	if !f.Matches(&ALBLog{Time: "2024-01-15T12:01:00.000000Z"}) {
		t.Error("expected time at before boundary to match")
	}
}

func TestFilter_Combined(t *testing.T) {
	m, _ := parseStatusPattern("5xx")
	f := &Filter{
		ElbStatus: &m,
		Method:    "POST",
	}

	code502 := 502
	code200 := 200
	// Both match
	if !f.Matches(&ALBLog{ElbStatusCode: &code502, RequestMethod: "POST"}) {
		t.Error("expected both conditions to match")
	}
	// Only status matches
	if f.Matches(&ALBLog{ElbStatusCode: &code502, RequestMethod: "GET"}) {
		t.Error("expected AND logic to reject")
	}
	// Only method matches
	if f.Matches(&ALBLog{ElbStatusCode: &code200, RequestMethod: "POST"}) {
		t.Error("expected AND logic to reject")
	}
}

func TestFilter_AlbGenerated(t *testing.T) {
	f := &Filter{AlbGenerated: true}

	code200 := 200
	// No target response (nil TargetStatusCode) → ALB-generated, should match
	if !f.Matches(&ALBLog{}) {
		t.Error("expected nil TargetStatusCode to match --alb-generated")
	}
	// Has target response → not ALB-generated, should be filtered out
	if f.Matches(&ALBLog{TargetStatusCode: &code200}) {
		t.Error("expected non-nil TargetStatusCode to be rejected by --alb-generated")
	}
}

func TestFilter_NoConditions(t *testing.T) {
	f := &Filter{}
	if !f.Matches(&ALBLog{}) {
		t.Error("empty filter should match everything")
	}
}

func TestParseTimeInput(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"2024-01-15T12:00:00Z", true},
		{"2024-01-15T12:00:00.123456Z", true},
		{"2024-01-15", true},
		{"not-a-time", false},
	}
	for _, tt := range tests {
		_, err := parseTimeInput(tt.input)
		if tt.valid && err != nil {
			t.Errorf("parseTimeInput(%q) unexpected error: %v", tt.input, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("parseTimeInput(%q) expected error", tt.input)
		}
	}
}
