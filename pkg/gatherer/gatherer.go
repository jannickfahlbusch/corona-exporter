package gatherer

import (
	"encoding/json"
	"log"
	"net/http"
	"pkg.jf-projects.de/corona-exporter/pkg/metrics"
	"time"
)

const DataURL = "https://www.bing.com/covid/data"

type retriever struct {
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

func Gather() {
	instance := &retriever{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for {
		log.Println("Gathering Metrics")
		entry, err := instance.Retrieve()
		if err != nil {
			log.Println(err)
		}

		setMetrics(entry)

		time.Sleep(1 * time.Hour)
	}
}

func setMetrics(entry *Entry) {
	for _, area := range entry.Areas {
		metrics.CasesPerCountry.WithLabelValues(area.DisplayName, string(metrics.CaseTypeConfirmed)).Set(float64(area.TotalConfirmed))
		metrics.CasesPerCountry.WithLabelValues(area.DisplayName, string(metrics.CaseTypeRecovered)).Set(float64(area.TotalRecovered))
		metrics.CasesPerCountry.WithLabelValues(area.DisplayName, string(metrics.CaseTypeCeased)).Set(float64(area.TotalDeaths))
	}

	metrics.CasesTotal.WithLabelValues(string(metrics.CaseTypeConfirmed)).Set(float64(entry.TotalConfirmed))
	metrics.CasesTotal.WithLabelValues(string(metrics.CaseTypeRecovered)).Set(float64(entry.TotalRecovered))
	metrics.CasesTotal.WithLabelValues(string(metrics.CaseTypeCeased)).Set(float64(entry.TotalDeaths))
}

func (instance *retriever) Retrieve() (*Entry, error) {
	response, err := instance.Client.Get(DataURL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	entry := &Entry{}
	err = json.NewDecoder(response.Body).Decode(entry)

	return entry, err
}
