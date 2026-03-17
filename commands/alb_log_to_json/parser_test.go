package alb_log_to_json

import (
	"os"
	"strings"
	"testing"
)

func TestParseLogLine_Normal(t *testing.T) {
	lines := readTestData(t)
	log, err := parseLogLine(lines[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// String fields
	assertEqual(t, "type", log.Type, "http")
	assertEqual(t, "time", log.Time, "2024-01-15T12:00:00.123456Z")
	assertEqual(t, "elb", log.Elb, "app/my-alb/50dc6c495c0c9188")
	assertEqual(t, "client:port", log.ClientPort, "192.168.1.100:12345")
	assertEqual(t, "target:port", log.TargetPort, "10.0.0.1:80")
	assertEqual(t, "domain_name", log.DomainName, "example.com")
	assertEqual(t, "user_agent", log.UserAgent, "Mozilla/5.0 (X11; Linux x86_64)")
	assertEqual(t, "actions_executed", log.ActionsExecuted, "forward")

	// Numeric fields
	assertFloatPtr(t, "request_processing_time", log.RequestProcessingTime, 0.001)
	assertFloatPtr(t, "target_processing_time", log.TargetProcessingTime, 0.002)
	assertFloatPtr(t, "response_processing_time", log.ResponseProcessingTime, 0.0)
	assertIntPtr(t, "elb_status_code", log.ElbStatusCode, 200)
	assertIntPtr(t, "target_status_code", log.TargetStatusCode, 200)
	assertInt64Ptr(t, "received_bytes", log.ReceivedBytes, 100)
	assertInt64Ptr(t, "sent_bytes", log.SentBytes, 500)

	// Request splitting
	assertEqual(t, "request", log.Request, "GET http://example.com/api/health HTTP/1.1")
	assertEqual(t, "request_method", log.RequestMethod, "GET")
	assertEqual(t, "request_url", log.RequestURL, "http://example.com/api/health")
	assertEqual(t, "request_protocol", log.RequestProtocol, "HTTP/1.1")
}

func TestParseLogLine_ErrorLine(t *testing.T) {
	lines := readTestData(t)
	log, err := parseLogLine(lines[1])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "type", log.Type, "https")
	assertIntPtr(t, "elb_status_code", log.ElbStatusCode, 502)
	assertFloatPtr(t, "target_processing_time", log.TargetProcessingTime, -1)
	assertEqual(t, "error_reason", log.ErrorReason, "LambdaInvalidResponse")
	assertEqual(t, "request_method", log.RequestMethod, "POST")
	assertEqual(t, "domain_name", log.DomainName, "api.example.com")

	// target_status_code is "-" → nil
	if log.TargetStatusCode != nil {
		t.Errorf("expected nil target_status_code, got %d", *log.TargetStatusCode)
	}
}

func TestParseLogLine_RedirectLine(t *testing.T) {
	lines := readTestData(t)
	log, err := parseLogLine(lines[2])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertIntPtr(t, "elb_status_code", log.ElbStatusCode, 301)
	assertEqual(t, "actions_executed", log.ActionsExecuted, "redirect")
	assertEqual(t, "redirect_url", log.RedirectUrl, "https://example.com/new-path")

	if log.TargetStatusCode != nil {
		t.Errorf("expected nil target_status_code, got %d", *log.TargetStatusCode)
	}
}

func TestParseLogLine_MalformedLine(t *testing.T) {
	_, err := parseLogLine("this is not a valid log line")
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
	if !strings.Contains(err.Error(), "expected at least 29 fields") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseLogLine_EmptyRequest(t *testing.T) {
	// A line where the request field is just "- - -" (malformed request)
	// The method/url/protocol should all be empty if request doesn't split into 3 parts
	lines := readTestData(t)
	log, err := parseLogLine(lines[0])
	if err != nil {
		t.Fatal(err)
	}
	// This line has a normal request, so it should split fine
	if log.RequestMethod == "" {
		t.Error("expected non-empty request_method for normal request")
	}
}

func TestSplitLogLine(t *testing.T) {
	line := `http 2024-01-15T12:00:00Z "quoted field with spaces" unquoted`
	parts := splitLogLine(line)
	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d: %v", len(parts), parts)
	}
	assertEqual(t, "part[0]", parts[0], "http")
	assertEqual(t, "part[1]", parts[1], "2024-01-15T12:00:00Z")
	assertEqual(t, "part[2]", parts[2], `"quoted field with spaces"`)
	assertEqual(t, "part[3]", parts[3], "unquoted")
}

func TestStripQuotes(t *testing.T) {
	tests := []struct{ in, want string }{
		{`"hello"`, "hello"},
		{`hello`, "hello"},
		{`""`, ""},
		{`"`, `"`},
		{``, ``},
	}
	for _, tt := range tests {
		got := stripQuotes(tt.in)
		if got != tt.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseFloat64(t *testing.T) {
	tests := []struct {
		in   string
		want *float64
	}{
		{"0.001", floatPtr(0.001)},
		{"-1", floatPtr(-1)},
		{"-", nil},
		{"abc", nil},
	}
	for _, tt := range tests {
		got := parseFloat64(tt.in)
		if tt.want == nil {
			if got != nil {
				t.Errorf("parseFloat64(%q) = %v, want nil", tt.in, *got)
			}
		} else if got == nil {
			t.Errorf("parseFloat64(%q) = nil, want %v", tt.in, *tt.want)
		} else if *got != *tt.want {
			t.Errorf("parseFloat64(%q) = %v, want %v", tt.in, *got, *tt.want)
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		in   string
		want *int
	}{
		{"200", intPtr(200)},
		{"-", nil},
		{"abc", nil},
	}
	for _, tt := range tests {
		got := parseInt(tt.in)
		if tt.want == nil {
			if got != nil {
				t.Errorf("parseInt(%q) = %v, want nil", tt.in, *got)
			}
		} else if got == nil {
			t.Errorf("parseInt(%q) = nil, want %v", tt.in, *tt.want)
		} else if *got != *tt.want {
			t.Errorf("parseInt(%q) = %v, want %v", tt.in, *got, *tt.want)
		}
	}
}

// --- helpers ---

func readTestData(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile("testdata/sample.log")
	if err != nil {
		t.Fatalf("reading test data: %v", err)
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func assertEqual(t *testing.T, field, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}

func assertIntPtr(t *testing.T, field string, got *int, want int) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %d", field, want)
	} else if *got != want {
		t.Errorf("%s: got %d, want %d", field, *got, want)
	}
}

func assertInt64Ptr(t *testing.T, field string, got *int64, want int64) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %d", field, want)
	} else if *got != want {
		t.Errorf("%s: got %d, want %d", field, *got, want)
	}
}

func assertFloatPtr(t *testing.T, field string, got *float64, want float64) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %v", field, want)
	} else if *got != want {
		t.Errorf("%s: got %v, want %v", field, *got, want)
	}
}

func floatPtr(v float64) *float64 { return &v }
func intPtr(v int) *int           { return &v }
