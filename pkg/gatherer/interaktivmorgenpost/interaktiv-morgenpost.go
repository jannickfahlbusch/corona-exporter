package interaktivmorgenpost

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"pkg.jf-projects.de/corona-exporter/pkg/gatherer"
	"pkg.jf-projects.de/corona-exporter/pkg/metrics"
)

// This is taken from https://github.com/montanaflynn/covid-19 which parses this nicely as JSON
const (
	DataURL = "https://raw.githubusercontent.com/montanaflynn/covid-19/master/data/current.json"
	Source  = "InteraktivMorgenpost"
)

type InterakivMorgenpost struct {
	client *http.Client
}

type cases struct {
	Updated   int `json:"updated"`
	Confirmed int `json:"confirmed"`
	Recovered int `json:"recovered"`
	Deaths    int `json:"deaths"`
}

type results struct {
	Global  map[string]cases `json:"global"`
	USA     map[string]cases `json:"usa"`
	Canada  map[string]cases `json:"canada"`
	Germany map[string]cases `json:"germany"`
	China   map[string]cases `json:"china"`
}

func NewInteraktivMorgenpost(client *http.Client) gatherer.Gatherer {
	return &InterakivMorgenpost{
		client: client,
	}
}

func (interaktivMorgenpost *InterakivMorgenpost) Gather() {
	for {
		log.Println("Gathering Metrics from Interaktiv Morgenpost")
		results, err := interaktivMorgenpost.Retrieve()
		if err != nil {
			log.Println(err)
			continue
		}

		interaktivMorgenpost.setCases(results)

		time.Sleep(1 * time.Hour)
	}
}

func (interaktivMorgenpost *InterakivMorgenpost) Retrieve() (*results, error) {
	response, err := interaktivMorgenpost.client.Get(DataURL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	results := &results{}
	err = json.NewDecoder(response.Body).Decode(results)

	return results, err
}

func (interaktivMorgenpost *InterakivMorgenpost) setCases(result *results) {
	var (
		totalConfirmed, totalRecovered, totalDeath int
	)

	for name, entry := range result.Global {
		setMetricsForArea(name, entry)
		totalConfirmed += entry.Confirmed
		totalRecovered += entry.Recovered
		totalDeath += entry.Deaths
	}

	for name, entry := range result.Germany {
		setMetricsForArea(name, entry)
	}

	for name, entry := range result.Canada {
		setMetricsForArea(name, entry)
	}

	for name, entry := range result.USA {
		setMetricsForArea(name, entry)
	}

	for name, entry := range result.China {
		setMetricsForArea(name, entry)
	}

	metrics.CasesTotal.WithLabelValues(Source, string(metrics.CaseTypeConfirmed)).Set(float64(totalConfirmed))
	metrics.CasesTotal.WithLabelValues(Source, string(metrics.CaseTypeRecovered)).Set(float64(totalRecovered))
	metrics.CasesTotal.WithLabelValues(Source, string(metrics.CaseTypeCeased)).Set(float64(totalDeath))
}

func setMetricsForArea(name string, entry cases) {
	metrics.CasesPerCountry.WithLabelValues(Source, name, string(metrics.CaseTypeConfirmed)).Set(float64(entry.Confirmed))
	metrics.CasesPerCountry.WithLabelValues(Source, name, string(metrics.CaseTypeRecovered)).Set(float64(entry.Recovered))
	metrics.CasesPerCountry.WithLabelValues(Source, name, string(metrics.CaseTypeCeased)).Set(float64(entry.Deaths))
}
