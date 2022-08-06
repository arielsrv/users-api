package common

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
	"sync"
)

var (
	padlock  = &sync.Mutex{}
	instance *MetricCollector
)

type MetricCollector struct {
	nrapp *newrelic.Application
}

func (r *MetricCollector) Nrapp() *newrelic.Application {
	return r.nrapp
}

func (r *MetricCollector) Record(name string, value float64) {
	r.nrapp.RecordCustomMetric(name, value)
}

func GetMetricCollector() *MetricCollector {
	if instance == nil {
		padlock.Lock()
		defer padlock.Unlock()
		if instance == nil {
			nrapp, _ := newrelic.NewApplication(
				newrelic.ConfigAppName("golang-users-api"),
				newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
				newrelic.ConfigDebugLogger(os.Stdout),
			)
			instance = &MetricCollector{nrapp: nrapp}
		}
	}

	return instance
}
