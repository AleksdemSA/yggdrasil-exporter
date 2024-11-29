package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var getpeers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "getpeers",
	Help: "get all peers",
})

var getrunnedpeers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "getrunnedpeers",
	Help: "get running peers",
})

var getsessions = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "getsessions",
	Help: "get sessions",
})

func executeCommandAndUpdateMetric(cmdStr string, metric prometheus.Gauge) {
	cmd := exec.Command("bash", "-c", cmdStr)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing command '%s': %v\n", cmdStr, err)
		return
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		log.Printf("Error parsing output '%s': %v\n", string(output), err)
		return
	}

	metric.Set(value)
}

func main() {
	prometheus.MustRegister(getpeers)
	prometheus.MustRegister(getrunnedpeers)
	prometheus.MustRegister(getsessions)

	go func() {
		for {
			executeCommandAndUpdateMetric("yggdrasilctl getpeers | wc -l", getpeers)
			executeCommandAndUpdateMetric("yggdrasilctl getpeers | grep Up | wc -l", getrunnedpeers)
			executeCommandAndUpdateMetric("yggdrasilctl getsessions | tail -n +2 | wc -l", getsessions)
			
			time.Sleep(60 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting exporter on :9120")
	log.Fatal(http.ListenAndServe(":9120", nil))
}

