package alb_log_to_json

import (
	"fmt"
	"strconv"
	"strings"
)

// parseLogLine parses a single ALB access log line into an ALBLog struct.
func parseLogLine(line string) (*ALBLog, error) {
	parts := splitLogLine(line)
	if len(parts) < 29 {
		return nil, fmt.Errorf("expected at least 29 fields, got %d", len(parts))
	}

	for i := range parts {
		parts[i] = stripQuotes(parts[i])
	}

	log := &ALBLog{
		Type:                 parts[0],
		Time:                 parts[1],
		Elb:                  parts[2],
		ClientPort:           parts[3],
		TargetPort:           parts[4],
		RequestProcessingTime:  parseFloat64(parts[5]),
		TargetProcessingTime:   parseFloat64(parts[6]),
		ResponseProcessingTime: parseFloat64(parts[7]),
		ElbStatusCode:          parseInt(parts[8]),
		TargetStatusCode:       parseInt(parts[9]),
		ReceivedBytes:          parseInt64(parts[10]),
		SentBytes:              parseInt64(parts[11]),
		Request:              parts[12],
		UserAgent:            parts[13],
		SslCipher:            parts[14],
		SslProtocol:          parts[15],
		TargetGroupArn:       parts[16],
		TraceId:              parts[17],
		DomainName:           parts[18],
		ChosenCertArn:        parts[19],
		MatchedRulePriority:  parts[20],
		RequestCreationTime:  parts[21],
		ActionsExecuted:      parts[22],
		RedirectUrl:          parts[23],
		ErrorReason:          parts[24],
		TargetPortList:       parts[25],
		TargetStatusCodeList: parts[26],
		Classification:       parts[27],
		ClassificationReason: parts[28],
	}

	// Split request into method, URL, protocol
	reqParts := strings.SplitN(log.Request, " ", 3)
	if len(reqParts) == 3 {
		log.RequestMethod = reqParts[0]
		log.RequestURL = reqParts[1]
		log.RequestProtocol = reqParts[2]
	}

	return log, nil
}

// splitLogLine splits an ALB log line by spaces, respecting quoted fields.
func splitLogLine(line string) []string {
	var parts []string
	var part strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case c == '"':
			inQuotes = !inQuotes
			part.WriteByte(c)
		case c == ' ' && !inQuotes:
			parts = append(parts, part.String())
			part.Reset()
		default:
			part.WriteByte(c)
		}
	}

	if part.Len() > 0 {
		parts = append(parts, part.String())
	}

	return parts
}

func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func parseFloat64(s string) *float64 {
	if s == "-" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

func parseInt(s string) *int {
	if s == "-" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func parseInt64(s string) *int64 {
	if s == "-" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}
