package main

import (
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"os"
)
var (
	appVersion = os.Getenv("IMAGE_TAG")
	backColor = "DarkGreen"
	imageBuildDate = os.Getenv("IMAGE_BUILD_DATE")
	kubeNodeName = os.Getenv("KUBE_NODE_NAME")
	kubePodName = os.Getenv("KUBE_POD_NAME")
	kubePodIP = os.Getenv("KUBE_POD_IP")
	hostname, _ = os.Hostname()
)

var (
	PrometheusHTTPRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "smackapi",
			Name: "http_request_count",
			Help: "The number of HTTP requests.",
		},
		[]string{"code", "appVersion", "backColor", "hostname", "request"},
	)

	PrometheusHTTPRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "smackapi",
			Name: "http_request_latency",
			Help: "The latency of HTTP requests.",
		},
		[]string{"code", "appVersion", "backColor", "hostname", "request"},
	)
)

func main() {

	if len(appVersion) == 0 {
		appVersion = "master-testing"
	}

	router := NewRouter()

	prometheus.MustRegister(PrometheusHTTPRequestCount)
	prometheus.MustRegister(PrometheusHTTPRequestLatency)

	log.Fatal(http.ListenAndServe(":8020", router))
}
