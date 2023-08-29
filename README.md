# Capital.com Prometheus Exporter

This repository contains the source code for a Prometheus exporter for the Capital.com broker. The exporter exposes various metrics related to the accounts such as balance, deposit, profit/loss, and availability.

## Table of Contents

- [Project Description](#project-description)
- [Features](#features)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Project Description

This exporter is designed to communicate with the Capital.com API and transform the response into metrics that Prometheus can scrape. The exporter uses the Go programming language and leverages the client_golang library from Prometheus.

## Features

- Fetches account details and balance info from the Capital.com API
- Exposes metrics to be scraped by Prometheus
- Uses the Go.uber.org/zap library for structured, leveled logging
- Configuration can be passed via environment variables or flags
- Supports both demo and real accounts

## Quick Start

You should have a working Go environment setup. If not, please follow [Install Go](https://golang.org/doc/install).

Clone the repository:

```bash
git clone https://github.com/<your-username>/capital-exporter.git
cd capital-exporter
```

## Build the program:

```bash
go build -o capital-exporter
```
## Environment Variables

```bash
export EXPORTER_email=<your-email>
export EXPORTER_password=<your-password>
export EXPORTER_apikey=<your-api-key>
export EXPORTER_port=<port-to-listen-on>  # Default is 9682
export EXPORTER_demo=<true-or-false>  # Default is false
export EXPORTER_interval=<interval-in-seconds>  # Default is 60
export EXPORTER_debug=<true-or-false>  # Default is false
```
## Usage

```bash
./capital-exporter
```

## Flags

You can also pass the configuration parameters as flags. Here is a list of all the flags:



- --email: Your Capital.com email
- --password: Your Capital.com password
- --apikey: Your Capital.com API Key
- --port: Port to listen on
- --demo: Use demo account
- --interval: Interval to update metrics
- --debug: Debug mode

Now you can see the metrics at `http://localhost:<EXPORTER_port>/metrics`

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

# License
This project is licensed under the MIT License. See the LICENSE file for details.

# Disclaimer

## This project is not affiliated with, maintained, authorized, endorsed, or sponsored by Capital.com or any of its affiliates. Use at your own risk.