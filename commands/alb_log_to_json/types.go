package alb_log_to_json

// ALBLog represents a parsed ALB access log entry with proper types.
type ALBLog struct {
	Type                   string   `json:"type"`
	Time                   string   `json:"time"`
	Elb                    string   `json:"elb"`
	ClientPort             string   `json:"client:port"`
	TargetPort             string   `json:"target:port"`
	RequestProcessingTime  *float64 `json:"request_processing_time"`
	TargetProcessingTime   *float64 `json:"target_processing_time"`
	ResponseProcessingTime *float64 `json:"response_processing_time"`
	ElbStatusCode          *int     `json:"elb_status_code"`
	TargetStatusCode       *int     `json:"target_status_code"`
	ReceivedBytes          *int64   `json:"received_bytes"`
	SentBytes              *int64   `json:"sent_bytes"`
	Request                string   `json:"request"`
	RequestMethod          string   `json:"request_method"`
	RequestURL             string   `json:"request_url"`
	RequestProtocol        string   `json:"request_protocol"`
	UserAgent              string   `json:"user_agent"`
	SslCipher              string   `json:"ssl_cipher"`
	SslProtocol            string   `json:"ssl_protocol"`
	TargetGroupArn         string   `json:"target_group_arn"`
	TraceId                string   `json:"trace_id"`
	DomainName             string   `json:"domain_name"`
	ChosenCertArn          string   `json:"chosen_cert_arn"`
	MatchedRulePriority    string   `json:"matched_rule_priority"`
	RequestCreationTime    string   `json:"request_creation_time"`
	ActionsExecuted        string   `json:"actions_executed"`
	RedirectUrl            string   `json:"redirect_url"`
	ErrorReason            string   `json:"error_reason"`
	TargetPortList         string   `json:"target:port_list"`
	TargetStatusCodeList   string   `json:"target_status_code_list"`
	Classification         string   `json:"classification"`
	ClassificationReason   string   `json:"classification_reason"`
}

// FieldRegistry maps JSON field names to accessor functions for --fields support.
var FieldRegistry = map[string]func(*ALBLog) interface{}{
	"type":                     func(l *ALBLog) interface{} { return l.Type },
	"time":                     func(l *ALBLog) interface{} { return l.Time },
	"elb":                      func(l *ALBLog) interface{} { return l.Elb },
	"client:port":              func(l *ALBLog) interface{} { return l.ClientPort },
	"target:port":              func(l *ALBLog) interface{} { return l.TargetPort },
	"request_processing_time":  func(l *ALBLog) interface{} { return l.RequestProcessingTime },
	"target_processing_time":   func(l *ALBLog) interface{} { return l.TargetProcessingTime },
	"response_processing_time": func(l *ALBLog) interface{} { return l.ResponseProcessingTime },
	"elb_status_code":          func(l *ALBLog) interface{} { return l.ElbStatusCode },
	"target_status_code":       func(l *ALBLog) interface{} { return l.TargetStatusCode },
	"received_bytes":           func(l *ALBLog) interface{} { return l.ReceivedBytes },
	"sent_bytes":               func(l *ALBLog) interface{} { return l.SentBytes },
	"request":                  func(l *ALBLog) interface{} { return l.Request },
	"request_method":           func(l *ALBLog) interface{} { return l.RequestMethod },
	"request_url":              func(l *ALBLog) interface{} { return l.RequestURL },
	"request_protocol":         func(l *ALBLog) interface{} { return l.RequestProtocol },
	"user_agent":               func(l *ALBLog) interface{} { return l.UserAgent },
	"ssl_cipher":               func(l *ALBLog) interface{} { return l.SslCipher },
	"ssl_protocol":             func(l *ALBLog) interface{} { return l.SslProtocol },
	"target_group_arn":         func(l *ALBLog) interface{} { return l.TargetGroupArn },
	"trace_id":                 func(l *ALBLog) interface{} { return l.TraceId },
	"domain_name":              func(l *ALBLog) interface{} { return l.DomainName },
	"chosen_cert_arn":          func(l *ALBLog) interface{} { return l.ChosenCertArn },
	"matched_rule_priority":    func(l *ALBLog) interface{} { return l.MatchedRulePriority },
	"request_creation_time":    func(l *ALBLog) interface{} { return l.RequestCreationTime },
	"actions_executed":         func(l *ALBLog) interface{} { return l.ActionsExecuted },
	"redirect_url":             func(l *ALBLog) interface{} { return l.RedirectUrl },
	"error_reason":             func(l *ALBLog) interface{} { return l.ErrorReason },
	"target:port_list":         func(l *ALBLog) interface{} { return l.TargetPortList },
	"target_status_code_list":  func(l *ALBLog) interface{} { return l.TargetStatusCodeList },
	"classification":           func(l *ALBLog) interface{} { return l.Classification },
	"classification_reason":    func(l *ALBLog) interface{} { return l.ClassificationReason },
}
