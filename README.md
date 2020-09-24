# SAP HA alert handlers


## Table of Content

- [Architecture](#architecture)
- [API](#api)
- [Quickstart](#quickstart)
- [Devel](#devel)


# Rationale

The alert handler serve as central component for building reactive self-healing systems.

It is main functionality is to handle Prometheus alerts automatically, which are fired by Prometheus.


# Architecture

See [design](doc/design.md)


# Quickstart

See [quickstart](doc/quickstart.md)

# API:

See [api](doc/api.md)

# Devel 

* build: `go build ./`

Deploy to the node

## Debugging:

`systemctl  restart  prometheus-alertmanager`

`prometheus-alertmanager`

`amtool` cli for alertmanager:
examples: https://github.com/prometheus/alertmanager#examples


## Debugging alerts

`promtool check rules /path/to/example.rules.yml`
