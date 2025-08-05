#!/usr/bin/env bash

CONTAINER_ENGINE=$(command -v podman 2>/dev/null || command -v docker 2>/dev/null)

$CONTAINER_ENGINE run -p 4318:4318 --rm -v "$(pwd)/collector-config.yaml:/etc/otelcol/config.yaml" otel/opentelemetry-collector-contrib --config /etc/otelcol/config.yaml
