## Clash Exporter

This is an exporter for Clash, for used by the [Prometheus](https://prometheus.io/) to monitor clash network traffic.

![](./images/grafana.png)

### Usage

#### shell
```sh
export CLASH_HOST="http://"
export CLASH_TOKEN="clash"

./clash-exporter -h                                                                                                                                                                                                  (base)
 
 Usage of ./clash-exporter:
  -collectDest
        enable collector dest
        Warning: collector destination if enabled, will generate a large number of metrics, which may put a lot of pressure on Prometheus.
  -port int
        port to listen on (default 2112)
```

#### use by docker
```sh
docker run -d --name clash-exporter -p 2112:2112 -e CLASH_HOST="${CLASH_HOST}" -e CLASH_TOKEN="$CLASH_TOKEN" zxh326/clash-exporter:latest 
```

#### use by systemctl
```sh
cat > /etc/systemd/system/clash-exporter.service << EOF
[Unit]
Description=Clash exporter
After=network.target

[Service]
Type=simple
Restart=always
Environment="CLASH_HOST=127.0.0.1:9090"
Environment="CLASH_TOKEN=clash"
ExecStart=/root/clash-exporter -collectDest

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable clash-exporter
systemctl start clash-exporter
```

####

visit http://localhost:2112/metrics and configure prometheus to scrape this endpoint.


### Prometheus Example Config
```yaml
  - job_name: "clash"
    metrics_path: /metrics
    scrape_interval: 1s
    static_configs:
      - targets: ["127.0.0.1:2112"]
```

### Grafana Example Dashboard

You can import [clash-dashboard.json](./grafana/dashboard.json) to obtain the example effect, or you can create one yourself based on the following metrics introduction.


### Metrics
| Metric name                       | Metric type | Labels                                                              |
|-----------------------------------|-------------|---------------------------------------------------------------------|
| clash_download_bytes_total        | Gauge       |                                                                     |
| clash_active_connections          | Gauge       |                                                                     |
| clash_network_traffic_bytes_total | Counter     | `soruce`,`destination(if enabled)`,`policy`,`type(download,upload)` |


### TODO
- [ ] dns query metrics
- [ ] proxy dial metrics