package alb_log_to_json

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	flagElbStatus    string
	flagTargetStatus string
	flagDomain       string
	flagMethod       string
	flagPath         string
	flagTargetGroup  string
	flagAfter        string
	flagBefore       string
	flagLimit        int
	flagFields       string
)

var Cmd = &cobra.Command{
	Use:   "alb-log-to-json [flags] [files...]",
	Short: "Convert ALB access logs to JSON",
	Long:  `Convert ALB access logs to JSON. Reads from stdin or specified files. Gzip files (.gz) are decompressed automatically.`,
	Args:  cobra.ArbitraryArgs,
	RunE:  runCmd,
}

func init() {
	Cmd.Flags().StringVar(&flagElbStatus, "elb-status", "", "Filter by ELB status code (e.g., 502, 5xx)")
	Cmd.Flags().StringVar(&flagTargetStatus, "target-status", "", "Filter by target status code (e.g., 200, 2xx)")
	Cmd.Flags().StringVar(&flagDomain, "domain", "", "Filter by domain name (substring match)")
	Cmd.Flags().StringVar(&flagMethod, "method", "", "Filter by request method (GET, POST, etc.)")
	Cmd.Flags().StringVar(&flagPath, "path", "", "Filter by request URL (substring match)")
	Cmd.Flags().StringVar(&flagTargetGroup, "target-group", "", "Filter by target group ARN (substring match)")
	Cmd.Flags().StringVar(&flagAfter, "after", "", "Filter entries at or after this time (RFC3339 or YYYY-MM-DD)")
	Cmd.Flags().StringVar(&flagBefore, "before", "", "Filter entries at or before this time (RFC3339 or YYYY-MM-DD)")
	Cmd.Flags().IntVar(&flagLimit, "limit", 0, "Maximum number of matching records to output (0 = unlimited)")
	Cmd.Flags().StringVar(&flagFields, "fields", "", "Comma-separated list of fields to include in output")
}

func runCmd(cmd *cobra.Command, args []string) error {
	filter, err := buildFilter()
	if err != nil {
		return err
	}

	var fields []string
	if flagFields != "" {
		for _, f := range strings.Split(flagFields, ",") {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			if _, ok := FieldRegistry[f]; !ok {
				return fmt.Errorf("unknown field: %q", f)
			}
			fields = append(fields, f)
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	count := 0

	if len(args) == 0 {
		return processReader(os.Stdin, encoder, filter, fields, &count)
	}

	for _, path := range args {
		if err := processFile(path, encoder, filter, fields, &count); err != nil {
			return err
		}
		if flagLimit > 0 && count >= flagLimit {
			break
		}
	}
	return nil
}

func buildFilter() (*Filter, error) {
	f := &Filter{}

	if flagElbStatus != "" {
		m, err := parseStatusPattern(flagElbStatus)
		if err != nil {
			return nil, err
		}
		f.ElbStatus = &m
	}

	if flagTargetStatus != "" {
		m, err := parseStatusPattern(flagTargetStatus)
		if err != nil {
			return nil, err
		}
		f.TargetStatus = &m
	}

	f.Domain = flagDomain
	f.Method = flagMethod
	f.Path = flagPath
	f.TargetGroup = flagTargetGroup

	if flagAfter != "" {
		t, err := parseTimeInput(flagAfter)
		if err != nil {
			return nil, fmt.Errorf("--after: %w", err)
		}
		f.After = &t
	}

	if flagBefore != "" {
		t, err := parseTimeInput(flagBefore)
		if err != nil {
			return nil, fmt.Errorf("--before: %w", err)
		}
		f.Before = &t
	}

	return f, nil
}

func processFile(path string, encoder *json.Encoder, filter *Filter, fields []string, count *int) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	var reader io.Reader = f
	if filepath.Ext(path) == ".gz" {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("reading gzip %s: %w", path, err)
		}
		defer gz.Close()
		reader = gz
	}

	return processReader(reader, encoder, filter, fields, count)
}

func processReader(r io.Reader, encoder *json.Encoder, filter *Filter, fields []string, count *int) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	for scanner.Scan() {
		if flagLimit > 0 && *count >= flagLimit {
			return nil
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		log, err := parseLogLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping malformed line: %v\n", err)
			continue
		}

		if !filter.Matches(log) {
			continue
		}

		if fields != nil {
			m := make(map[string]interface{}, len(fields))
			for _, f := range fields {
				m[f] = FieldRegistry[f](log)
			}
			if err := encoder.Encode(m); err != nil {
				return err
			}
		} else {
			if err := encoder.Encode(log); err != nil {
				return err
			}
		}

		*count++
	}

	return scanner.Err()
}
