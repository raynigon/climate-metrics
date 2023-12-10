package config

import (
	"github.com/raynigon/climate-metrics/v2/collector"
)

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func (cfg ClimateMetricsExporterConfig) GetCollectorConfig() collector.CollectorConfig {
	return collector.CollectorConfig{
		Logger: cfg.GetLogger(),
		GPIO:   *cfg.dhtGPIO,
	}
}
