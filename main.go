package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zxh326/clash-exporter/collector"
)

var collectDest bool
var port int

func init() {
	flag.BoolVar(&collectDest, "collectDest", false, "enable collector dest\nWarning: if collector destination enabled, will generate a large number of metrics, which may put a lot of pressure on Prometheus.")
	flag.IntVar(&port, "port", 2112, "port to listen on")
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func main() {
	flag.Parse()

	config := collector.CollectConfig{
		CollectDest: collectDest,
		ClashHost:   getEnvOrDefault("CLASH_HOST", "127.0.0.1:9090"),
		ClashToken:  getEnvOrDefault("CLASH_TOKEN", ""),
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Listening on :", port)

	go collector.Start(config)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
