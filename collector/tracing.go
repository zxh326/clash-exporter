package collector

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	dnsReqDuration    *prometheus.HistogramVec
	ruleMatchDuration *prometheus.HistogramVec
	proxyDialDuration *prometheus.HistogramVec
)

type tracingMessage struct {
	Duration int      `json:"duration"`
	ID       string   `json:"id"`
	Metadata Metadata `json:"metadata"`
	Payload  string   `json:"payload"`
	Proxy    string   `json:"proxy"`
	Rule     string   `json:"rule"`
	Type     string   `json:"type"`
	DnsType  string   `json:"dnsType"`
}

type Tracing struct {
}

func (t *Tracing) Name() string {
	return "tracing"
}

func (t *Tracing) Collect(config CollectConfig) error {
	if !config.CollectTracing {
		return nil
	}
	log.Println("starting collector:", t.Name())
	ctx := context.Background()
	endpoint := fmt.Sprintf("ws://%s/profile/tracing", config.ClashHost)
	if config.ClashToken != "" {
		endpoint = fmt.Sprintf("%s?token=%s", endpoint, config.ClashToken)
	}
	conn, resp, err := websocket.Dial(ctx, endpoint, nil)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			log.Fatal("profile tracing is not enabled in Clash, please enable it first, or use `-collectTracing=false` disable this collector. \nFYI: https://github.com/Dreamacro/clash/wiki/Clash-Premium-Features#tracing")
		}
		log.Fatal("tracing: failed to dial: ", err)
	}

	conn.SetReadLimit(1024 * 1024)

	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	for {
		var m tracingMessage
		err = wsjson.Read(ctx, conn, &m)
		if err != nil {
			return errors.Wrap(err, "failed to read JSON message")
		}
		switch m.Type {
		case "RuleMatch":
			ruleMatchDuration.WithLabelValues().Observe(float64(m.Duration) / 1000)
		case "DNSRequest":
			dnsReqDuration.WithLabelValues(m.DnsType).Observe(float64(m.Duration) / 1000)
		case "ProxyDial":
			proxyDialDuration.WithLabelValues(m.Proxy).Observe(float64(m.Duration) / 1000)
		}
	}

}

func init() {
	dnsReqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "clash",
			Subsystem: "tracing",
			Name:      "dns_request_duration_milliseconds",
			Help:      "DNS request duration in milliseconds",
			Buckets:   []float64{1, 5, 10, 25, 50, 75, 100, 250, 500}, // 1ms ~ 500ms
		},
		[]string{"type"},
	)

	ruleMatchDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "clash",
			Subsystem: "tracing",
			Name:      "rule_match_duration_milliseconds",
			Help:      "Rule match duration in milliseconds",
			Buckets:   []float64{1, 5, 10, 25, 50, 75, 100, 250, 500}, // 1ms ~ 500ms
		},
		[]string{},
	)

	proxyDialDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "clash",
			Subsystem: "tracing",
			Name:      "proxy_dial_duration_milliseconds",
			Help:      "Proxy dial duration in milliseconds",
			Buckets:   []float64{1, 10, 50, 100, 250, 500, 1000, 2500, 5000}, // 1ms ~ 5s
		},
		[]string{"policy"},
	)

	prometheus.MustRegister(dnsReqDuration, ruleMatchDuration, proxyDialDuration)

	t := &Tracing{}
	Register(t)
}
