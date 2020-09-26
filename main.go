package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	handlerWebSrvPort  = "9999"
	sapHandlerName     = "/hooks-sap"
	defautlHandlerName = "/hooks-default"
)

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
	h := new(HanaDiskFull)
	// make sure we run only 1 handler until it finish.
	http.HandleFunc(sapHandlerName, h.handlerHanaDiskFull)

	http.HandleFunc(defautlHandlerName, defaultHandler)

	// TODO: serve webserver (future https)
	err := http.ListenAndServe(":"+handlerWebSrvPort, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
