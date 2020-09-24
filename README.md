# SAP HA alert handlers


## Table of Content

- [Architecture](doc/design.md)
- [API](doc/api.md)
- [Quickstart](doc/quickstart.md)
- [Devel](#devel)


# Rationale

The alert handler serve as central component for building reactive self-healing systems, based on Prometheus.

It is meant to be a single binary, which will be distribued and installed to all nodes of clusters.
The routing itself is delegated to alertmanager of prometheus.
It is main functionality is to `selfheal` and handle Prometheus alerts automatically, which are fired by Prometheus.

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
