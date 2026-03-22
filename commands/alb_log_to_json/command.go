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
	flagAlbGenerated bool
	flagCount        bool
	flagGroupBy      string
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
	Cmd.Flags().BoolVar(&flagAlbGenerated, "alb-generated", false, "Filter to ALB-generated errors (no target response)")
	Cmd.Flags().BoolVar(&flagCount, "count", false, "Count matching records instead of outputting JSON")
	Cmd.Flags().StringVar(&flagGroupBy, "group-by", "", "Group and count by field (from FieldRegistry)")
}

func runCmd(cmd *cobra.Command, args []string) error {
	filter, err := buildFilter()
	if err != nil {
		return err
	}

	// Validate mutual exclusivity
	if flagGroupBy != "" && flagCount {
		return fmt.Errorf("--group-by and --count are mutually exclusive")
	}
	if flagGroupBy != "" && flagFields != "" {
		return fmt.Errorf("--group-by and --fields are mutually exclusive")
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

	// Set up aggregator for --group-by
	var aggregator *Aggregator
	if flagGroupBy != "" {
		if _, ok := FieldRegistry[flagGroupBy]; !ok {
			return fmt.Errorf("unknown field for --group-by: %q", flagGroupBy)
		}
		aggregator = NewAggregator(flagGroupBy)
	}

	var encoder *json.Encoder
	if !flagCount && aggregator == nil {
		encoder = json.NewEncoder(os.Stdout)
	}

	count := 0

	if len(args) == 0 {
		if err := processReader(os.Stdin, encoder, filter, fields, aggregator, &count); err != nil {
			return err
		}
	} else {
		for _, path := range args {
			if err := processFile(path, encoder, filter, fields, aggregator, &count); err != nil {
				return err
			}
			if flagLimit > 0 && count >= flagLimit {
				break
			}
		}
	}

	if flagCount {
		fmt.Println(count)
	}
	if aggregator != nil {
		enc := json.NewEncoder(os.Stdout)
		for _, r := range aggregator.Results() {
			if err := enc.Encode(r); err != nil {
				return err
			}
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
	f.AlbGenerated = flagAlbGenerated

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

func processFile(path string, encoder *json.Encoder, filter *Filter, fields []string, aggregator *Aggregator, count *int) error {
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

	return processReader(reader, encoder, filter, fields, aggregator, count)
}

func processReader(r io.Reader, encoder *json.Encoder, filter *Filter, fields []string, aggregator *Aggregator, count *int) error {
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

		if flagCount || aggregator != nil {
			if aggregator != nil {
				aggregator.Add(log)
			}
			*count++
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
