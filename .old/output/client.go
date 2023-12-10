package output

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

type OutputClient struct {
	influxClient influxdb2.Client
	database     string
	measurement  string
	tags         map[string]string
}

func NewClient(url string, username string, password string, database string, measurement string, tags map[string]string) *OutputClient {
	var client = OutputClient{
		influxClient: influxdb2.NewClient(url, username+":"+password),
		database:     database,
		measurement:  measurement,
		tags:         tags,
	}
	return &client
}

func (c *OutputClient) WritePoints(humidity float64, temperature float64) {
	writeApi := c.influxClient.WriteApiBlocking("", c.database+"/autogen")
	p := influxdb2.NewPoint("room_climate",
		c.tags,
		map[string]interface{}{"humidity": humidity, "temperature": temperature},
		time.Now())
	writeApi.WritePoint(context.Background(), p)
}

func (c *OutputClient) Close() {
	c.influxClient.Close()
}
