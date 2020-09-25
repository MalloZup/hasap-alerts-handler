package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
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

type CrmMon struct {
	XMLName        xml.Name `xml:"crm_mon"`
	Version        string   `xml:"version,attr"`
	NodeAttributes struct {
		Node []struct {
			Name      string `xml:"name,attr"`
			Attribute []struct {
				Name  string `xml:"name,attr"`
				Value string `xml:"value,attr"`
			} `xml:"attribute"`
		} `xml:"node"`
	} `xml:"node_attributes"`
}

func isHanaNodePrimary() (bool, error) {
	var c *CrmMon
	crmMonXML, err := exec.Command("/usr/sbin/crm_mon", "-X", "--inactive").Output()
	if err != nil {
		return false, fmt.Errorf("error while executing crm_mon: %w", err)
	}

	err = xml.Unmarshal(crmMonXML, &c)
	if err != nil {
		return false, fmt.Errorf("error while parsing crm_mon XML %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return false, fmt.Errorf("error could not get hostname %w", err)
	}

	log.Debugf("isHanaNodePrimary method called, hostname: %s", hostname)

	// TODO: verify if we can rely safely on this assumption
	// that the node name of cluster(CIB) is equal to hostname

	for _, n := range c.NodeAttributes.Node {
		if n.Name != hostname {
			continue
		}
		// check if primary attr is set, then hana is primary
		for _, a := range n.Attribute {
			//  <attribute name="hana_prd_site" value="PRIMARY_SITE_NAME"/>
			if a.Value == "PRIMARY_SITE_NAME" {
				return true, nil
			}
		}
	}

	return false, nil
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
	log.Infof("HanaDiskFullHandler called by %s", alerts.Receiver)

	// iterate over alerts since we could have a group of them
	for _, a := range alerts.Alerts {
		log.Infof("%s generated by %s", a.Labels.Alertname, a.GeneratorURL)

		// we look only for hana components
		if strings.ToLower(a.Labels.Component) != "hana" {
			continue
		}
		// check if self-healing is enabled otherwise skip
		if strings.ToLower(a.Labels.Selfhealing) != "true" {
			log.Debugln("selfhealing disabled")
			continue
		}
		primary, err := isHanaNodePrimary()
		if err != nil {
			log.Warnf("[CRITICAL] Error by detecting if hana is primary node %s", err)
		}

		if primary == true {
			// todo: get the resource name dinamically, eg. the postfix
			hanaResource := "msl_SAPHana_PRD_HDB00"
			cmd := exec.Command("/usr/sbin/crm", "resource", "move", hanaResource, "force")
			log.Infoln("[SELFHEALING]: selfhealing HANA primary node. Migrating to other node")
			err := cmd.Run()
			if err != nil {
				log.Warnln("[CRITICAL]: move resource hana resource")
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

	// register the various handler
	http.HandleFunc("/hooks-sap", handlerHanaDiskFull)
	http.HandleFunc("/hooks-default", defaultHandler)

	// TODO: serve webserver (future https)
	err := http.ListenAndServe(":"+handlerWebSrvPort, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
