# cloud-utils

A CLI for working with cloud services and tools.

## Install

```zsh
go install github.com/igor-kupczynski/cloud-utils@latest
```

## alb-log-to-json

Convert AWS ALB access logs to JSON. Reads from stdin or files (`.gz` auto-decompressed), with built-in filtering and field selection.

```zsh
# Basic: pipe from stdin
cat alb_access_log.txt | cloud-utils alb-log-to-json

# Read files directly (gzip supported)
cloud-utils alb-log-to-json /tmp/alb-logs/*.gz

# Filter 5xx errors
cloud-utils alb-log-to-json --elb-status 5xx /tmp/alb-logs/*.gz

# Filter by domain and method, limit output
cloud-utils alb-log-to-json --domain api.example.com --method POST --limit 50 *.gz

# Time range
cloud-utils alb-log-to-json --after 2024-01-15 --before 2024-01-16 *.gz

# Select specific fields to reduce output
cloud-utils alb-log-to-json --fields time,elb_status_code,request_url,target_processing_time *.gz

# Combine with jq for further analysis (numeric types work natively)
cloud-utils alb-log-to-json --elb-status 5xx --fields time,target_processing_time,request_url *.gz \
  | jq -s 'sort_by(.target_processing_time) | reverse | .[:10]'
```

For the full list of flags, run:

```zsh
cloud-utils alb-log-to-json --help
```

## Build from source

```zsh
git clone https://github.com/igor-kupczynski/cloud-utils.git
cd cloud-utils
make build   # build locally
make install # install to $GOPATH/bin
make test    # run tests
```