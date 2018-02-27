package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type Config struct {
	Key          string `json:"Key"`
	BackColor    string `json:"BackColor"`
	AppVersion   string `json:"AppVersion"`
	BuildDate    string `json:"BuildDate"`
	KubeNodeName string `json:"KubeNodeName"`
	KubePodName  string `json:"KubePodName"`
	KubePodIP    string `json:"KubePodIP"`
}

type Configs []Config

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RUNNING")
}

func returnConfig(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	code := http.StatusInternalServerError

	defer func() { // Make sure we record a status.
		duration := time.Since(start)
		PrometheusHTTPRequestCount.WithLabelValues(fmt.Sprintf("%v", code), appVersion, backColor, hostname, "getconfig").Inc()
		PrometheusHTTPRequestLatency.WithLabelValues(fmt.Sprintf("%v", code), appVersion, backColor, hostname, "getconfig").Observe(duration.Seconds())
	}()

	configs := Config{Key: "10", BackColor: backColor, AppVersion: appVersion, BuildDate: imageBuildDate, KubeNodeName: kubeNodeName, KubePodName: kubePodName, KubePodIP: kubePodIP}

	// insert simulated delay if color is red
	var delay int
	if strings.Contains(strings.ToUpper(backColor), "RED") {
		//delay = random(50, 1000)
		delay = 500
	} else if strings.Contains(strings.ToUpper(backColor), "BLUE"){
		//delay = random(20, 200)
		delay = 100
	} else {
		delay = 50
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(configs); err != nil {
		panic(err)
	}

	code = http.StatusOK
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func testHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "RUNNING")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check. can simulate error with http status
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusBadGateway)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func metricHandler(w http.ResponseWriter, r *http.Request) {
	prometheus.Handler().ServeHTTP(w, r)
}