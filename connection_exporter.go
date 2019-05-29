package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	//"os"
	"time"
	"os/exec"
	"strconv"
	"strings"
)

var (
	connection_timewait = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "connection_timewait",
		Help: "connection timewait count",
	})
	connection_established = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "connection_established",
		Help: "connection established count",
	})
)

func init() {
	prometheus.MustRegister(connection_timewait)
	prometheus.MustRegister(connection_established)
}

func connectionStatus(metric string) int {
	cmd :=exec.Command("/bin/ss","-antp")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("%s", err)
	}
	result := strings.Count(string(out),metric)
	return result
}
func recordMetrics() {
	go func() {
		for {
			timewait := connectionStatus("TIME-WAIT")
			established := connectionStatus("ESTAB")
			log.Info("Time Wait connections: "+ strconv.Itoa(timewait))
			log.Info("Established connections: "+ strconv.Itoa(established))
			connection_established.Set(float64(established))
			connection_timewait.Set(float64(timewait))
			time.Sleep(10 * time.Second)
		}
	}()
}

func main() {
	bind := flag.String("bind", "0.0.0.0", "bind port default 0.0.0.0")
	port := flag.Int("port", 9319, "port to listen default 9319")
	flag.Parse()
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port " + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(*bind+":"+strconv.Itoa(*port), nil))
}
