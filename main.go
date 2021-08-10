package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

type config struct {
	timeout          time.Duration
	port             uint
	verbose          bool
	connectionString string
}

func readAndValidateConfig() config {
	var result config

	flag.StringVar(&result.connectionString, "connection-string", "", "Starboard connection string")
	flag.UintVar(&result.port, "port", 9580, "Port to expose scraping endpoint on")
	flag.DurationVar(&result.timeout, "timeout", time.Second*30, "Timeout for scrape")
	flag.BoolVar(&result.verbose, "verbose", false, "Enable verbose logging")

	flag.Parse()

	if result.connectionString == "" {
		log.Fatal("Starboard connection string not provided")
	}

	log.WithFields(logrus.Fields{
		"port":    result.port,
		"timeout": result.timeout,
		"verbose": result.verbose,
	}).Infof("Starboard exporter configured")

	return result
}

func configureRoutes() {
	var landingPage = []byte(`<html>
		<head><title>Starboard exporter for Prometheus</title></head>
		<body>
		<h1>Starboard exporter for Prometheus</h1>
		<p><a href='/metrics'>Metrics</a></p>
		</body>
		</html>
		`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage) // nolint: errcheck
	})

	http.Handle("/metrics", promhttp.Handler())
}

func setupLogger(config config) {
	if config.verbose {
		log.Level = logrus.DebugLevel
	}
}

func startHTTPServer(config config) {
	listenAddr := fmt.Sprintf(":%d", config.port)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func main() {

	config := readAndValidateConfig()
	setupLogger(config)

	configureRoutes()

	// client := sb.New(config.connectionString, config.timeout)
	// coll := collector.New(client, log)
	// prometheus.MustRegister(coll)

	startHTTPServer(config)
}
