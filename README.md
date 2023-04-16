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
        It must be the Clash premium version, and the profile.tracing must be enabled in the Clash configuration file. (default false)
  -port int
        port to listen on (default 2112)
```

#### use by docker

```sh
CLASH_HOST=127.0.0.1:9090
CLASH_TOKEN=pass
docker run -d --name clash-exporter -p 2112:2112 -e CLASH_HOST="${CLASH_HOST}" -e CLASH_TOKEN="$CLASH_TOKEN" zzde/clash-exporter:latest
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

- You can import [clash-dashboard.json](./grafana/dashboard.json) to obtain the example effect, or you can create one yourself based on the following metrics introduction.

- or Import via [grafana.com](https://grafana.com/grafana/dashboards/18530-clash-dashboard/) with id `18530`

### Metrics

| Metric name                                     | Metric type | Labels                                                              |
| ----------------------------------------------- | ----------- | ------------------------------------------------------------------- |
| clash_info                                      | Gauge       | `version`, `premium`                                                |
| clash_download_bytes_total                      | Gauge       |                                                                     |
| clash_upload_bytes_total                        | Gauge       |                                                                     |
| clash_active_connections                        | Gauge       |                                                                     |
| clash_network_traffic_bytes_total               | Counter     | `source`,`destination(if enabled)`,`policy`,`type(download,upload)` |
| clash_tracing_rule_match_duration_milliseconds  | Histogram   |                                                                     |
| clash_tracing_dns_request_duration_milliseconds | Histogram   | `type(dnsType)`                                                     |
| clash_tracing_proxy_dial_duration_milliseconds  | Histogram   | `policy`                                                            |

### FAQ

- tracing metrics is empty

  - Required clash premium version
  - Follow [clash profile docs](https://github.com/Dreamacro/clash/wiki/Clash-Premium-Features#tracing) enable profile tracing
  - Add `-collectTracing=true` flag in clash-exporter start script

- high Prometheus Memory

  This may be caused by the default enable of collector destination traffic, which can generate a large number of metrics. Try use `-collectDest=false` disable it.

### TODO

- [x] dns query metrics
- [x] proxy dial metrics
