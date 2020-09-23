package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const handlerWebSrvPort = "9999"

type PromAlert struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Summary string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		Severity  string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}

func handler1(_ http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var t PromAlert
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)

	log.Infoln("handler called via prometheus alert")
}

func main() {
	log.Infof("starting handler on port: %s", handlerWebSrvPort)
	// register the various handlers
	http.HandleFunc("/hooks-default", handler1)

	// serve webserver (future https)
	err := http.ListenAndServe(":"+handlerWebSrvPort, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
