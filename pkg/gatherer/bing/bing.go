package bing

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"pkg.jf-projects.de/corona-exporter/pkg/gatherer"
	"pkg.jf-projects.de/corona-exporter/pkg/metrics"
)

const (
	DataURL = "https://www.bing.com/covid/data"
	Source  = "Bing"
)

type Bing struct {
	Client *http.Client
}

type Entry struct {
	Areas []*Entry `json:"areas"`

	TotalConfirmed int64 `json:"totalConfirmed"`
	TotalDeaths    int64 `json:"totalDeaths"`
	TotalRecovered int64 `json:"totalRecovered"`

	ID          string    `json:"id"`
	LastUpdated time.Time `json:"lastUpdated"`
	DisplayName string    `json:"displayName"`
	Lat         float64   `json:"lat"`
	Long        float64   `json:"long"`
	Country     string    `json:"country"`
	ParentID    string    `json:"parentId"`
}

func NewBingGatherer(client *http.Client) gatherer.Gatherer {
	return &Bing{
		Client: client,
	}
}

func (bing *Bing) Gather() {
	for {
		log.Println("Gathering Metrics from Bing")
		entry, err := bing.Retrieve()
		if err != nil {
			log.Println(err)
			continue
		}

		setMetrics(entry)

		time.Sleep(1 * time.Hour)
	}
}

func setMetrics(entry *Entry) {
	for _, area := range entry.Areas {
		metrics.CasesPerCountry.WithLabelValues(Source, area.DisplayName, string(metrics.CaseTypeConfirmed)).Set(float64(area.TotalConfirmed))
		metrics.CasesPerCountry.WithLabelValues(Source, area.DisplayName, string(metrics.CaseTypeRecovered)).Set(float64(area.TotalRecovered))
		metrics.CasesPerCountry.WithLabelValues(Source, area.DisplayName, string(metrics.CaseTypeCeased)).Set(float64(area.TotalDeaths))
	}

	metrics.CasesTotal.WithLabelValues(Source, string(metrics.CaseTypeConfirmed)).Set(float64(entry.TotalConfirmed))
	metrics.CasesTotal.WithLabelValues(Source, string(metrics.CaseTypeRecovered)).Set(float64(entry.TotalRecovered))
	metrics.CasesTotal.WithLabelValues(Source, string(metrics.CaseTypeCeased)).Set(float64(entry.TotalDeaths))
}

func (bing *Bing) Retrieve() (*Entry, error) {
	request, err := http.NewRequest(http.MethodGet, DataURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept-Language", "en")
	response, err := bing.Client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	entry := &Entry{}
	err = json.NewDecoder(response.Body).Decode(entry)

	return entry, err
}
