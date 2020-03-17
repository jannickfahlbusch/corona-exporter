package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	Namespace = "corona"
	Cases     = "cases"
)

type CaseType string

var (
	CaseTypeConfirmed CaseType = "confirmed"
	CaseTypeRecovered CaseType = "recovered"
	CaseTypeCeased    CaseType = "ceased"
)

type GaugeVec interface {
}

var (
	CasesPerCountry = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Cases,
			Name:      "per_country",
			Help:      "Corona Cases by Country and type (Confirmed, Recovered, Ceased)",
		}, []string{"country", "type"})

	CasesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Cases,
			Name:      "total",
			Help:      "Corona Cases in total by type (Confirmed, Recovered, Ceased)",
		}, []string{"type"})
)

func init() {
	prometheus.MustRegister(CasesPerCountry, CasesTotal)
}
