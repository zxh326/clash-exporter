## Clash Exporter

This is an exporter for Clash, for used by the [Prometheus](https://prometheus.io/) to monitor clash network traffic.

![](./images/grafana.png)

### Usage

```
➜  ./clash-exporter -h
Usage of ./clash-exporter:
  -collectDest
        enable collector dest
        Warning: if collector destination enabled, will generate a large number of metrics, which may put a lot of pressure on Prometheus. (default true)
  -collectTracing
        enable collector tracing.
        It must be the Clash premium version, and the profile.tracing must be enabled in the Clash configuration file. (default true)
  -port int
        port to listen on (default 2112)
```

#### use by docker

```sh
docker run -d --name clash-exporter -p 2112:2112 -e CLASH_HOST="${CLASH_HOST}" -e CLASH_TOKEN="$CLASH_TOKEN" ghcr.io/zxh326/clash-exporter:latest
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

| Metric name                                     | Metric type | Labels                                                              |
| ----------------------------------------------- | ----------- | ------------------------------------------------------------------- |
| clash_download_bytes_total                      | Gauge       |                                                                     |
| clash_upload_bytes_total                        | Gauge       |                                                                     |
| clash_active_connections                        | Gauge       |                                                                     |
| clash_network_traffic_bytes_total               | Counter     | `soruce`,`destination(if enabled)`,`policy`,`type(download,upload)` |
| clash_tracing_rule_match_duration_milliseconds  | Histogram   |                                                                     |
| clash_tracing_dns_request_duration_milliseconds | Histogram   | `type(dnsType)`                                                     |
| clash_tracing_proxy_dial_duration_milliseconds  | Histogram   |                                                                     |

### TODO

- [x] dns query metrics
- [x] proxy dial metrics
