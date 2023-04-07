package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Info struct{}

var (
	clashInfo *prometheus.GaugeVec
)

func (*Info) Name() string {
	return "info"
}

type versionMessage struct {
	Version string `json:"version"`
	Premium bool   `json:"premium"`
}

func (*Info) Collect(config CollectConfig) error {
	endpoint := fmt.Sprintf("http://%s/version", config.ClashHost)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	if config.ClashToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.ClashToken))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("version: failed to read response body: ", err)
	}
	if responseBody == nil {
		return nil
	}

	var result versionMessage
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return err
	}

	clashInfo.WithLabelValues(result.Version, fmt.Sprintf("%t", result.Premium)).Set(1)

	return nil
}

func init() {
	clashInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "clash",
			Name:      "info",
			Help:      "Clash Infos",
		},
		[]string{"version", "premium"},
	)
	prometheus.MustRegister(clashInfo)
	Register(new(Info))
}
