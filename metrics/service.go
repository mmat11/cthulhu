package metrics

import (
	kitmetrics "github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "cthulhu"

type Service interface {
	IncUpdatesTotal(groupName string, userName string)
	ObserveUpdatesDuration(groupName string, ms float64)
}

type metrics struct {
	UpdatesTotal           kitmetrics.Counter
	UpdatesDurationSeconds kitmetrics.Histogram
}

func NewService() Service {
	return &metrics{
		UpdatesTotal: kitprometheus.NewCounterFrom(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "updates_total",
				Help:      "total number of updates labeled by group name",
			},
			[]string{"group", "user"},
		),
		UpdatesDurationSeconds: kitprometheus.NewHistogramFrom(
			prometheus.HistogramOpts{
				Namespace: Namespace,
				Name:      "updates_duration_seconds",
				Help:      "duration of updates labeled by group name",
			},
			[]string{"group"},
		),
	}
}

func (m *metrics) IncUpdatesTotal(groupName string, userName string) {
	m.UpdatesTotal.With("group", groupName).With("user", userName).Add(1)
}

func (m *metrics) ObserveUpdatesDuration(groupName string, seconds float64) {
	m.UpdatesDurationSeconds.With("group", groupName).Observe(seconds)
}
