package alb_log_to_json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type ALBLog struct {
	Type                   string `json:"type"`
	Time                   string `json:"time"`
	Elb                    string `json:"elb"`
	ClientPort             string `json:"client:port"`
	TargetPort             string `json:"target:port"`
	RequestProcessingTime  string `json:"request_processing_time"`
	TargetProcessingTime   string `json:"target_processing_time"`
	ResponseProcessingTime string `json:"response_processing_time"`
	ElbStatusCode          string `json:"elb_status_code"`
	TargetStatusCode       string `json:"target_status_code"`
	ReceivedBytes          string `json:"received_bytes"`
	SentBytes              string `json:"sent_bytes"`
	Request                string `json:"request"`
	UserAgent              string `json:"user_agent"`
	SslCipher              string `json:"ssl_cipher"`
	SslProtocol            string `json:"ssl_protocol"`
	TargetGroupArn         string `json:"target_group_arn"`
	TraceId                string `json:"trace_id"`
	DomainName             string `json:"domain_name"`
	ChosenCertArn          string `json:"chosen_cert_arn"`
	MatchedRulePriority    string `json:"matched_rule_priority"`
	RequestCreationTime    string `json:"request_creation_time"`
	ActionsExecuted        string `json:"actions_executed"`
	RedirectUrl            string `json:"redirect_url"`
	ErrorReason            string `json:"error_reason"`
	TargetPortList         string `json:"target:port_list"`
	TargetStatusCodeList   string `json:"target_status_code_list"`
	Classification         string `json:"classification"`
	ClassificationReason   string `json:"classification_reason"`
}

var Cmd = &cobra.Command{
	Use:   "alb-log-to-json",
	Short: "Convert ALB access logs to JSON",
	Long:  `Convert ALB access logs to JSON by reading from stdin`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			logLine := scanner.Text()
			logParts := []string{}
			part := ""

			for i := 0; i < len(logLine); i++ {
				c := logLine[i]

				if c == ' ' && strings.Count(part, "\"")%2 == 0 {
					logParts = append(logParts, part)
					part = ""
				} else {
					part += string(c)
				}
			}

			if part != "" {
				logParts = append(logParts, part)
			}

			for i, part := range logParts {
				if i == 17 || i == 18 || i == 19 {
					if len(part) >= 2 && part[0] == '"' && part[len(part)-1] == '"' {
						logParts[i] = part[1 : len(part)-1]
					}
				}
			}

			albLog := &ALBLog{
				Type:                   logParts[0],
				Time:                   logParts[1],
				Elb:                    logParts[2],
				ClientPort:             logParts[3],
				TargetPort:             logParts[4],
				RequestProcessingTime:  logParts[5],
				TargetProcessingTime:   logParts[6],
				ResponseProcessingTime: logParts[7],
				ElbStatusCode:          logParts[8],
				TargetStatusCode:       logParts[9],
				ReceivedBytes:          logParts[10],
				SentBytes:              logParts[11],
				Request:                logParts[12],
				UserAgent:              logParts[13],
				SslCipher:              logParts[14],
				SslProtocol:            logParts[15],
				TargetGroupArn:         logParts[16],
				TraceId:                logParts[17],
				DomainName:             logParts[18],
				ChosenCertArn:          logParts[19],
				MatchedRulePriority:    logParts[20],
				RequestCreationTime:    logParts[21],
				ActionsExecuted:        logParts[22],
				RedirectUrl:            logParts[23],
				ErrorReason:            logParts[24],
				TargetPortList:         logParts[25],
				TargetStatusCodeList:   logParts[26],
				Classification:         logParts[27],
				ClassificationReason:   logParts[28],
			}

			albLog.Type = strings.Trim(albLog.Type, "\"")
			albLog.Request = strings.Trim(albLog.Request, "\"")
			albLog.UserAgent = strings.Trim(albLog.UserAgent, "\"")
			albLog.ActionsExecuted = strings.Trim(albLog.ActionsExecuted, "\"")
			albLog.RedirectUrl = strings.Trim(albLog.RedirectUrl, "\"")
			albLog.ErrorReason = strings.Trim(albLog.ErrorReason, "\"")
			albLog.TargetPortList = strings.Trim(albLog.TargetPortList, "\"")
			albLog.TargetStatusCodeList = strings.Trim(albLog.TargetStatusCodeList, "\"")
			albLog.Classification = strings.Trim(albLog.Classification, "\"")
			albLog.ClassificationReason = strings.Trim(albLog.ClassificationReason, "\"")

			jsonData, err := json.MarshalIndent(albLog, "", "  ")
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Println(string(jsonData))
		}
	},
}

func init() {
}
