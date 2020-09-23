# SAP HA alert handlers

# Rationale

The alert handler serve as central component for building reactive self-healing systems.

It is main functionality is to handle Prometheus alerts automatically, which are fired by Prometheus.


# Architecture Overview

See [doc/design.md](doc/design.md)


# Developer 

* build: `go build ./`

Deploy to the node

# Debugging:

`systemctl  restart  prometheus-alertmanager`

`prometheus-alertmanager`

`amtool` cli for alertmanager:
examples: https://github.com/prometheus/alertmanager#examples

