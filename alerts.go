package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// PromAlert is the alert payload alertmanager sent to the hook
type PromAlert struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname   string `json:"alertname"`
			Selfhealing string `json:"selfhealing`
			Severity    string `json:"severity"`
			Component   string `json:"component"`
			Instance    string `json:"instance"`
		} `json:"labels"`
		Annotations struct {
			Summary string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
}

// AlertFire is the payload from https://prometheus.io/docs/alerting/latest/clients/
type AlertFire struct {
	Status string `json:"status"`
	Labels struct {
		Alertname string `json:"alertname"`
		Severity  string `json:"severity"`
		Component string `json:"component"`
		Instance  string `json:"instance"`
	} `json:"labels"`
	Annotations struct {
		Summary string `json:"summary"`
	} `json:"annotations"`
	GeneratorURL string `json:"generatorURL"`
}

func (alert *AlertFire) sendAlert(url string) {
	body := &alert
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Errorf("Error sending http post alert %s", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)
	// Print the body to the stdout
	io.Copy(os.Stdout, res.Body)
}
