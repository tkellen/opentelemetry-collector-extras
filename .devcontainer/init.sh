#!/bin/sh
# OpenTelemetry Collector Builder
go install go.opentelemetry.io/collector/cmd/builder@latest
# OpenTelemetry Collector Metadata Generator
go install go.opentelemetry.io/collector/cmd/mdatagen@latest

# Go remote debugger
go install github.com/go-delve/delve/cmd/dlv@latest
# Go local debugger
go install github.com/google/gops@latest
