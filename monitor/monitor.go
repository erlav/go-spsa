package monitor

import (
	"fmt"
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PushMonitor struct {
	Pusher     *push.Pusher
	Keys_label []string
	Labels     map[string]string
}

func getRawURL() (string, error) {
	var scheme, host, port string
	var e bool = false

	port, e = os.LookupEnv("PG_PORT")

	if host, e = os.LookupEnv("PG_HOST"); e {
		if addr, err := url.Parse(host); err == nil {
			if scheme = addr.Scheme; scheme == "" {
				scheme = "http"
			}
			if port == "" {
				port = addr.Port()
			}
			host = addr.Host
		}
	}
	return fmt.Sprintf("%s://%s:%s", scheme, host, port), nil
}

func New(job string, keys_label []string, data map[string]string) (*PushMonitor, error) {
	rawurl, _ := getRawURL()
	var pm PushMonitor = PushMonitor{push.New(rawurl, job), keys_label, data}
	pm.Pusher.Grouping("tipointegracion", "api")
	for k, v := range data {
		pm.Pusher.Grouping(k, v)
	}
	return &pm, nil
}

func (pm *PushMonitor) PushMetric(name string, value float64, labels map[string]string) error {
	var g prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
	})
	g.Set(value)
	for k, v := range labels {
		pm.Pusher.Grouping(k, v)
	}

	return pm.Pusher.Collector(g).Push()
}
