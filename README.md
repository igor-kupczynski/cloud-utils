# cloud-utils

A CLI for working with cloud services and tools.

## Features

- Convert ALB access logs to JSON (`alb-log-to-json` command)

## Install

```zsh
go install github.com/igor-kupczynski/cloud-utils@latest
```

## Usage

To use the `alb-log-to-json` command, pipe the ALB access log to the `cloud-utils` binary:

```zsh
cat alb_access_log.txt | ./cloud-utils alb-log-to-json
```

## Build

1. Clone the repository:
```zsh
git clone https://github.com/igor-kupczynski/cloud-utils.git
cd cloud-utils
```

2. Build the project:
```zsh
make build
```