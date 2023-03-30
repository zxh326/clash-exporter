package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type message struct {
	DownloadTotal int64         `json:"downloadTotal"`
	UploadTotal   int64         `json:"uploadTotal"`
	Connections   []Connections `json:"connections"`
}
type Metadata struct {
	Network         string `json:"network"`
	Type            string `json:"type"`
	SourceIP        string `json:"sourceIP"`
	DestinationIP   string `json:"destinationIP"`
	SourcePort      string `json:"sourcePort"`
	DestinationPort string `json:"destinationPort"`
	Host            string `json:"host"`
	DNSMode         string `json:"dnsMode"`
	ProcessPath     string `json:"processPath"`
	SpecialProxy    string `json:"specialProxy"`
}
type Connections struct {
	ID          string    `json:"id"`
	Metadata    Metadata  `json:"metadata"`
	Upload      int       `json:"upload"`
	Download    int       `json:"download"`
	Start       time.Time `json:"start"`
	Chains      []string  `json:"chains"`
	Rule        string    `json:"rule"`
	RulePayload string    `json:"rulePayload"`
}

var (
	uploadTotalBytes   *prometheus.GaugeVec
	downloadTotalBytes *prometheus.GaugeVec
	activeConnections  *prometheus.GaugeVec

	todoDownloadTotal *prometheus.CounterVec
)

func handleConnections() {
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s%s?token=%s", CLASH_HOST, CONNECTIONS_PATH, CLASH_TOKEN), nil)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetReadLimit(1024 * 1024)

	connectionCache := make(map[string]float64)

	defer conn.Close(websocket.StatusInternalError, "the sky is falling")
	for {
		var m message
		err = wsjson.Read(ctx, conn, &m)
		if err != nil {
			log.Fatalf("wsjson.Read error: %v", err)
		}
		uploadTotalBytes.WithLabelValues().Set(float64(m.UploadTotal))
		downloadTotalBytes.WithLabelValues().Set(float64(m.DownloadTotal))
		activeConnections.WithLabelValues().Set(float64(len(m.Connections)))
		for _, connection := range m.Connections {
			if _, ok := connectionCache[connection.ID]; !ok {
				connectionCache[connection.ID] = 0
			}
			destination := connection.Metadata.Host
			if destination == "" {
				destination = connection.Metadata.DestinationIP
			}
			todoDownloadTotal.WithLabelValues(connection.Metadata.SourceIP, destination, connection.Chains[0]).Add(float64(connection.Download) - connectionCache[connection.ID])
			connectionCache[connection.ID] = float64(connection.Download)
		}
	}
}

func init() {
	uploadTotalBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "clash",
			Name:      "upload_total_bytes",
			Help:      "Total upload bytes",
		},
		[]string{},
	)
	downloadTotalBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "clash",
			Name:      "download_total_bytes",
			Help:      "Total download bytes",
		},
		[]string{},
	)

	activeConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "clash",
			Name:      "active_connections",
			Help:      "Active connections",
		},
		[]string{},
	)

	todoDownloadTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "clash",
			Name:      "host_rule_download_total",
			Help:      "Total download bytes by host rule",
		},
		[]string{"source", "destination", "policy"},
	)
	prometheus.MustRegister(uploadTotalBytes, downloadTotalBytes, todoDownloadTotal)
}
