# cloud-utils

A CLI for working with cloud services and tools.

## Features

- Convert ALB access logs to JSON (`alb-log-to-json` command)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/igor-kupczynski/cloud-utils.git
cd cloud-utils
```

2. Build the project:

```bash
make build
```

## Usage

To use the `alb-log-to-json` command, pipe the ALB access log to the `cloud-utils` binary:

```bash
cat alb_access_log.txt | ./cloud-utils alb-log-to-json
```