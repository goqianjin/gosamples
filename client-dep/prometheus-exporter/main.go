package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
)

var myCounter *prometheus.CounterVec

func main() {
	http.HandleFunc("/add", triggerMetricsHandler)

	register := prometheus.NewRegistry()
	myCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qj_demo_exporter_metrics_req_counter",
			Help: "qj_demo_exporter_metrics_req_counter",
		},
		[]string{"firstPath"},
	)
	register.MustRegister(myCounter)
	http.Handle("/metrics", promhttp.HandlerFor(register, promhttp.HandlerOpts{}))
	err := http.ListenAndServe(":10000", nil)
	log.Fatal("http.ListenAndServe:", err)
}

func triggerMetricsHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	pathSegements := strings.Split(urlPath, "/")
	if len(pathSegements) == 0 {
		fmt.Fprint(w, "Invalid path.")
	}
	myCounter.WithLabelValues(pathSegements[0]).Inc()
	fmt.Fprint(w, "added one metrics on: " + pathSegements[0])

}
