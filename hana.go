package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	//"github.com/MalloZup/hasap-alerts-handler/internal"
	log "github.com/sirupsen/logrus"
)

// CrmMon is the cluster pacemaker xml status via crm_mon
type CrmMon struct {
	XMLName   xml.Name `xml:"crm_mon"`
	Version   string   `xml:"version,attr"`
	Resources struct {
		Clone []struct {
			ID             string `xml:"id,attr"`
			MultiState     string `xml:"multi_state,attr"`
			Unique         string `xml:"unique,attr"`
			Managed        string `xml:"managed,attr"`
			Failed         string `xml:"failed,attr"`
			FailureIgnored string `xml:"failure_ignored,attr"`
			Resource       []struct {
				ID             string `xml:"id,attr"`
				ResourceAgent  string `xml:"resource_agent,attr"`
				Role           string `xml:"role,attr"`
				Active         string `xml:"active,attr"`
				Orphaned       string `xml:"orphaned,attr"`
				Blocked        string `xml:"blocked,attr"`
				Managed        string `xml:"managed,attr"`
				Failed         string `xml:"failed,attr"`
				FailureIgnored string `xml:"failure_ignored,attr"`
				NodesRunningOn string `xml:"nodes_running_on,attr"`
				Node           struct {
					Name   string `xml:"name,attr"`
					ID     string `xml:"id,attr"`
					Cached string `xml:"cached,attr"`
				} `xml:"node"`
			} `xml:"resource"`
		} `xml:"clone"`
	} `xml:"resources"`
}

// HanaDiskFull type
// Prometheus can send depending on interval same alerts multiple times in short interval
// of time. like 5s (depending on timeout)
// Prevent that an handler is called multiple times until it performing operation
// the mutex will ensure that we run only 1 handler until it finish
type HanaDiskFull struct {
	mu sync.Mutex
}

// handle when Hana Primary node has disk full
func (ns *HanaDiskFull) handlerHanaDiskFull(_ http.ResponseWriter, req *http.Request) {
	// see type description for this mutex
	ns.mu.Lock()
	defer ns.mu.Unlock()

	cmd := exec.Command("sleep", "30")
	log.Infoln("--- SLEEPING")
	cmd.Run()

	// read body json from Prometheus alertmanager
	decoder := json.NewDecoder(req.Body)
	var alerts PromAlert
	err := decoder.Decode(&alerts)
	if err != nil {
		log.Warnf("error by decoding json from hana-handler http: %s", err)
	}

	log.Infof("HanaDiskFullHandler called by %s", alerts.Receiver)

	for _, a := range alerts.Alerts {
		log.Infof("%s generated by %s", a.Labels.Alertname, a.GeneratorURL)

		// look only for hana components and selfhealing true labels
		if strings.ToLower(a.Labels.Component) != "hana" &&
			strings.ToLower(a.Labels.Selfhealing) != "true" {
			continue
		}

		// read crm_mon xml for detecting if hana is primary on node
		var cMon *CrmMon
		crmMonXML, err := exec.Command("/usr/sbin/crm_mon", "-X", "--inactive").Output()
		if err != nil {
			log.Warnf("error while executing crm_mon: %w", err)
			return
		}

		err = xml.Unmarshal(crmMonXML, &cMon)
		if err != nil {
			log.Warnf("error while parsing crm_mon XML %w", err)
			return
		}
		// used to see if the hana primary resource is running on local node
		nodeHostname, err := os.Hostname()
		if err != nil {
			log.Warnf("error could not get hostname %w", err)
			return
		}
		// check if the current node where the alert is executed
		// has the HANADB as primary db, then return the res name if yes
		// this will be used for taking over the res
		primaryRes := lookUpHanaNodePrimary(cMon, nodeHostname)

		if primaryRes != "" {
			cmd := exec.Command("/usr/sbin/crm", "resource", "move", primaryRes, "force")
			log.Infoln("[SELFHEALING]: selfhealing HANA primary node. Migrating to other node")
			err := cmd.Run()
			if err != nil {
				log.Warnln("[CRITICAL]: move resource hana resource")
			}

			// wait until the hana primary has migrated

		}
	}
}

// TODO: it could be that this is not safe since node could be an array TO verify it
func lookUpHanaNodePrimary(cmon *CrmMon, hostname string) string {
	for _, n := range cmon.Resources.Clone {
		matched, _ := regexp.MatchString(`msl_SAPHana_.*`, n.ID)
		if matched {
			for _, r := range n.Resource {
				if r.Role == "Master" && r.Node.Name == hostname {
					return n.ID
				}
			}
		}
	}
	return ""
}
