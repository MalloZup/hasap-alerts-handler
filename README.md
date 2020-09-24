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

The alert handlers needs a common API shared with the alertmanager. This is shared via labels.
Since labels names can be choosed freely, here is listed what is the convention for our alerts handler.

Labels which are influencing handler behaviour:

* `selfhealing:`: 

`Description:` Disable or not selfhealing(this disable the handler)
`Values allowed`: true/false or absent. If it is absent it is same as false.

* `component:` 
`Description:` This specify on which node/s the action of self-healing will run.
`Values allowed`: the values will be added.. WORKING IN PROGRESS (this will be update)


If the selfhealing is set true, the handler will selfheal based on the alert.

Example:
```
  - alert: AlertExampleSelfhealing
    expr: YourAlertExpr
    labels:
      severity: critical
      selfhealing: true 
      component: drbd # run on drbd nodes
    annotations:
      summary: drbd critical
```

If the labels selfhealing is false, or absent the handler will not work.

Example:
```
  - alert: AlertExampleNotHealing
    expr: YourAlertExpr
    labels:
      severity: critical
      selfhealing: false
      component: drbd
    annotations:
      summary: drbd critical
```

```
  - alert: AlertExampleNotHealing
    expr: YourAlertExpr
    labels:
      severity: critical
      selfhealing: false
      component: drbd
    annotations:
      summary: drbd critical
```

___
Labels which are not influencing the handler behaviour

**NOTE** In future, severity might influence how the alert-handler react to the alerts, as kind of     priority for scheduling. Choose the severity label accordingly.
*severity:  ( Critical,Major, Warning, Medium, Low) 
  

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
