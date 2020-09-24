# API

The alert handlers needs a common API shared with the alertmanager. This is shared via labels.
Since labels names can be choosed freely, here is listed what is the convention for our alerts handler.

## Labels which are influencing handler behaviour:

### `selfhealing:`: 

* `Description:` Disable or not selfhealing (this disable the handler)
If the selfhealing is set true, the handler will selfheal based on the alert.
* `Values allowed`: true/false or absent. If it is absent it is same as false.

### `component:` 
* `Description:` This specify on which node/s the action of self-healing will run.
* `Values allowed`: the values will be added.. WORKING IN PROGRESS (this will be update)

Examples:

This alert will be scheduled on drbd nodes and the handler will be called (selfhealing true)
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


The following labels run on drbd but no self-healing happens. See false and absent value.
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
## Labels which are not influencing the handler behaviour

**NOTE** In future, severity might influence how the alert-handler react to the alerts, as kind of     priority for scheduling. Choose the severity label accordingly.

### severity: 
* `Description`:  Describe the severity of an alert
* `Values`: Critical,Major,  Medium, Warning, Low
  
