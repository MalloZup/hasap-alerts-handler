package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// TODO handle this via config file later on
const handlerWebSrvPort = "9999"

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

// check if the node where the handler run, is the hana primary node
func isHanaNodePrimary() bool {
	// TODO: do something here
	return true
}

// handle when Hana Primary node has disk full
func handlerHanaDiskFull(_ http.ResponseWriter, req *http.Request) {
	// read body json from Prometheus alertmanager
	decoder := json.NewDecoder(req.Body)
	var alerts PromAlert
	err := decoder.Decode(&alerts)
	if err != nil {
		log.Warnln(err)
	}
	log.Infoln("HanaDiskFullHandler called")

	// iterate over alerts
	for _, a := range alerts.Alerts {
		// we look only for hana components
		if strings.ToLower(a.Labels.Component) != "hana" {
			continue
		}
		// check if self-healing is enabled otherwise skip
		if strings.ToLower(a.Labels.Selfhealing) != "true" {
			log.Infoln("selfhealing disabled")
			continue
		}

		// TODO: we need to know/check if the hana is primary node
		if isHanaNodePrimary() == true {
			cmd := exec.Command("sleep", "1")
			log.Infoln("[SELFHEALING]: selfhealing HANA primary node. Migrating to other node")
			err := cmd.Run()
			if err != nil {
				log.Warnln("[CRITICAL]: Could not selfhealing hana primary node")
			}
		}
	}
}

// default handler. this is where the alerts witch doesn't match anything goes
func defaultHandler(_ http.ResponseWriter, req *http.Request) {
	// read body json from Prometheus alertmanager
	decoder := json.NewDecoder(req.Body)
	var alerts PromAlert
	err := decoder.Decode(&alerts)
	if err != nil {
		log.Warnln(err)
	}
	// the default handler for moment does nothing
}

func main() {
	log.Infof("starting handler on port: %s", handlerWebSrvPort)
	// register the various handlers

	http.HandleFunc("/hooks-sap", handlerHanaDiskFull)

	http.HandleFunc("/hooks-default", defaultHandler)

	// TODO: serve webserver (future https)
	err := http.ListenAndServe(":"+handlerWebSrvPort, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
