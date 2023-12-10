package config

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type ClimateMetricsExporterConfig struct {
	listenAddress *string
	metricsPath   *string
	logLevel      *string
	logFormat     *string
	logOutput     *string
	dhtGPIO       *string
}

func NewClimateMetricsExporterConfig() ClimateMetricsExporterConfig {
	config := ClimateMetricsExporterConfig{
		listenAddress: kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").
			Default(":9776").
			Envar("CM_WEB_LISTEN_ADDRESS").
			String(),
		metricsPath: kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").
			Default("/metrics").
			Envar("CM_WEB_TELEMETRY_PATH").
			String(),
		logLevel: kingpin.Flag("log.level", "Sets the loglevel. Valid levels are debug, info, warn, error").
			Default("info").
			Envar("CM_LOG_LEVEL").
			String(),
		logFormat: kingpin.Flag("log.format", "Sets the log format. Valid formats are json and logfmt").
			Default("logfmt").
			Envar("CM_LOG_FORMAT").
			String(),
		logOutput: kingpin.Flag("log.output", "Sets the log output. Valid outputs are stdout and stderr").
			Default("stdout").
			Envar("CM_LOG_OUTPUT").
			String(),
		dhtGPIO: kingpin.Flag("dht.gpio", "The GPIO Pin used for the DHT sensor (e.g. GPIO2)").
			Default("GPIO2").
			Envar("CM_DHT_GPIO").
			String(),
	}
	kingpin.Version("0.0.1")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return config
}

func (cfg ClimateMetricsExporterConfig) GetListeningAccess() string {
	return *cfg.listenAddress
}

func (cfg ClimateMetricsExporterConfig) GetMetricsPath() string {
	return *cfg.metricsPath
}
