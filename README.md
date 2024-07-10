# go-autumn-restclient-circuitbreaker-prometheus

Prometheus instrumentation functions for 
[go-autumn-restclient-circuitbreaker](https://github.com/StephanHCB/go-autumn-restclient-circuitbreaker).

## About go-autumn-restclient

It's a rest client that also supports x-www-form-urlencoded.

## About go-autumn-restclient-circuitbreaker

This library adds another wrapper to the rest client which provides a circuit breaker.

## About go-autumn-restclient-circuitbreaker-prometheus

Implements instrumentation callbacks that use [prometheus/client_golang](https://github.com/prometheus/client_golang).

## Usage

Use the provided callbacks while constructing your rest client stack.
