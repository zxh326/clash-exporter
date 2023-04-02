package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dnsReqDuration *prometheus.HistogramVec
	dnsReqTotal    *prometheus.CounterVec
)

type Tracing struct {
}

func (t *Tracing) Name() string {
	return "tracing"
}

func (t *Tracing) Collect(config CollectConfig) error {
	// TODO: tracing
	return nil
}

func init() {
	dnsReqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "clash",
			Subsystem: "tracing",
			Name:      "dns_request_duration_seconds",
			Help:      "DNS request duration in seconds",
		},
		[]string{"domain", "type", "result"},
	)

	dnsReqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "clash",
			Subsystem: "tracing",
			Name:      "dns_request_total",
			Help:      "DNS request total",
		},
		[]string{"domain", "type", "result"},
	)

	prometheus.MustRegister(dnsReqTotal, dnsReqDuration)

	t := &Tracing{}
	register(t)
}
