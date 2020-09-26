package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlertFiring(t *testing.T) {
	var a *AlertFire
	a = new(AlertFire)
	a.Status = "firing"
	a.Labels.Alertname = "testing-alert"
	a.Labels.Component = "unit-test component"
	a.Labels.Severity = "critical"
	a.Labels.Instance = "test instance"
	a.Annotations.Summary = "just a test"
	a.GeneratorURL = "unit-test"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read body json from Prometheus alertmanager
		decoder := json.NewDecoder(r.Body)
		alert := make([]AlertFire, 1)
		decoder.Decode(&alert)
		t.Log(alert[0].Status)
		t.Log(alert[0].Labels.Component)
		if alert[0].GeneratorURL != "unit-test" {
			t.Errorf("got %s expected unit-test", alert[0].GeneratorURL)
		}
	}))
	defer ts.Close()

	a.sendAlert(ts.URL)
}
