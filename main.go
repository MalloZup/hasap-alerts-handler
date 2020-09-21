package main

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {

	handler1 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc handler called via Prometheus alert\n")
		log.Infoln("handler called via prometheus alert")
	}

	http.HandleFunc("/hooks", handler1)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
