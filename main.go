package main

import (
	"encoding/json"

	"fmt"
	"net/http"
	"os"

	"runtime"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	handlerWebSrvPort  = "9999"
	version            = "0.0.1"
	sapHandlerName     = "/hooks-sap"
	defautlHandlerName = "/hooks-default"
)

var (
	// the time the binary was built
	buildDate string
	// global --help flag
	helpFlag *bool
	// global --version flag
	versionFlag *bool
)

func showVersion() {
	fmt.Printf("version %s\nbuilt with %s %s/%s %s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH, buildDate)
	os.Exit(0)
}

func init() {
	flag.StringP("port", "p", handlerWebSrvPort, "The port number to listen on for http alerts event")
	flag.StringP("alertserver", "a", "", "prometheus alertmanager server IP")
	flag.StringP("config", "c", "", "The path to a custom configuration file. NOTE: it must be in yaml format.")
	flag.CommandLine.SortFlags = false

	helpFlag = flag.BoolP("help", "h", false, "show this help message")
	versionFlag = flag.BoolP("version", "v", false, "show version and build information")
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

	flag.Parse()

	switch {
	case *helpFlag:
		flag.Usage()
		os.Exit(0)
	case *versionFlag:
		fmt.Printf("version %s\nbuilt with %s %s/%s %s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH, buildDate)
		os.Exit(0)
	}

	configFile, err := Config(flag.CommandLine)
	if err != nil {
		log.Fatalf("Could not initialize config: %s", err)
	}

	log.Info(configFile.Get("port"))

	log.Infof("starting handler on port: %s", handlerWebSrvPort)

	// register the various handler
	h := new(HanaDiskFull)
	// make sure we run only 1 handler until it finish.
	http.HandleFunc(sapHandlerName, h.handlerHanaDiskFull)

	http.HandleFunc(defautlHandlerName, defaultHandler)

	// TODO: serve webserver (future https)
	err = http.ListenAndServe(":"+handlerWebSrvPort, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
