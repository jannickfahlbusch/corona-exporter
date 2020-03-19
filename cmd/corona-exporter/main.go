package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"pkg.jf-projects.de/corona-exporter/pkg/gatherer/bing"
	"pkg.jf-projects.de/corona-exporter/pkg/gatherer/interaktivmorgenpost"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()

	go startGatherers()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func startGatherers() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	go bing.NewBingGatherer(httpClient).Gather()
	go interaktivmorgenpost.NewInteraktivMorgenpost(httpClient).Gather()
}
