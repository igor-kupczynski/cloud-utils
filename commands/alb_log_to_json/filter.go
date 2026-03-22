package alb_log_to_json

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Filter holds parsed filter criteria. All conditions are AND-ed.
type Filter struct {
	ElbStatus    *statusMatcher
	TargetStatus *statusMatcher
	Domain       string
	Method       string
	Path         string
	TargetGroup  string
	After        *time.Time
	Before       *time.Time
	AlbGenerated bool
}

// Matches returns true if the log entry passes all filter conditions.
func (f *Filter) Matches(log *ALBLog) bool {
	if f.ElbStatus != nil && !f.ElbStatus.matches(log.ElbStatusCode) {
		return false
	}
	if f.TargetStatus != nil && !f.TargetStatus.matches(log.TargetStatusCode) {
		return false
	}
	if f.AlbGenerated && log.TargetStatusCode != nil {
		return false
	}
	if f.Domain != "" && !strings.Contains(log.DomainName, f.Domain) {
		return false
	}
	if f.Method != "" && !strings.EqualFold(log.RequestMethod, f.Method) {
		return false
	}
	if f.Path != "" && !strings.Contains(log.RequestURL, f.Path) {
		return false
	}
	if f.TargetGroup != "" && !strings.Contains(log.TargetGroupArn, f.TargetGroup) {
		return false
	}
	if f.After != nil || f.Before != nil {
		t, err := time.Parse(time.RFC3339Nano, log.Time)
		if err != nil {
			// Try without fractional seconds
			t, err = time.Parse(time.RFC3339, log.Time)
			if err != nil {
				return false
			}
		}
		if f.After != nil && t.Before(*f.After) {
			return false
		}
		if f.Before != nil && t.After(*f.Before) {
			return false
		}
	}
	return true
}

type statusMatcher struct {
	exact *int
	class int // e.g., 5 for 5xx
}

func (m *statusMatcher) matches(code *int) bool {
	if code == nil {
		return false
	}
	if m.exact != nil {
		return *code == *m.exact
	}
	return *code/100 == m.class
}

// parseStatusPattern parses a status pattern like "502" or "5xx".
func parseStatusPattern(pattern string) (statusMatcher, error) {
	if strings.HasSuffix(pattern, "xx") {
		prefix := strings.TrimSuffix(pattern, "xx")
		class, err := strconv.Atoi(prefix)
		if err != nil || class < 1 || class > 9 {
			return statusMatcher{}, fmt.Errorf("invalid status pattern: %s", pattern)
		}
		return statusMatcher{class: class}, nil
	}
	code, err := strconv.Atoi(pattern)
	if err != nil {
		return statusMatcher{}, fmt.Errorf("invalid status code: %s", pattern)
	}
	return statusMatcher{exact: &code}, nil
}

// parseTimeInput parses a time string in RFC3339 or YYYY-MM-DD format.
func parseTimeInput(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s (expected RFC3339 or YYYY-MM-DD)", s)
}
