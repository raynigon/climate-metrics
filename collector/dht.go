package collector

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/raynigon/climate-metrics/v2/pkg/dht"
)

var (
	dhtSubsystem = "dht"
)

type DhtCollector struct {
	config      CollectorConfig
	dht         *dht.DHT
	temperature *prometheus.Desc
	humidity    *prometheus.Desc
}

func init() {
	registerCollector(dhtSubsystem, NewDhtCollector)
}

// NewOrgActionsCollector returns a new Collector exposing actions billing stats.
func NewDhtCollector(config CollectorConfig, ctx context.Context) (Collector, error) {
	err := dht.HostInit()
	if err != nil {
		return nil, err
	}
	dht, err := dht.NewDHT(config.GPIO, dht.Fahrenheit, "11")
	if err != nil {
		return nil, err
	}
	collector := &DhtCollector{
		config: config,
		dht:    dht,
		temperature: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dhtSubsystem, "temperature"),
			"Temperature in Celsius",
			[]string{},
			nil,
		),
		humidity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dhtSubsystem, "humidity"),
			"Humidity in %",
			[]string{},
			nil,
		),
	}
	return collector, nil
}

// Describe implements Collector.
func (c *DhtCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperature
	ch <- c.humidity
}

func (c *DhtCollector) Reload(ctx context.Context) error {
	return nil
}

func (c *DhtCollector) Update(ctx context.Context, ch chan<- prometheus.Metric) error {
	temperature, humidity, err := c.dht.Read()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.temperature, prometheus.GaugeValue, float64(temperature))
	ch <- prometheus.MustNewConstMetric(c.humidity, prometheus.GaugeValue, float64(humidity))
	return nil
}
