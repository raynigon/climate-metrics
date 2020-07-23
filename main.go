package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/raynigon/climate-metrics/dht"
	"github.com/raynigon/climate-metrics/output"
	"gopkg.in/yaml.v2"
)
// Configuration contains all config values
type Configuration struct {
	Input struct {
		Gpio            string  `yaml:"gpio"`
		RefreshInterval float64 `yaml:"refreshInterval"`
	} `yaml:"input"`
	Output struct {
		URL         string            `yaml:"url"`
		Username    string            `yaml:"username"`
		Password    string            `yaml:"password"`
		Database    string            `yaml:"database"`
		Measurement string            `yaml:"measurement"`
		Tags        map[string]string `yaml:"tags"`
	} `yaml:"output"`
}

func readConfig() (Configuration, error) {
	var config Configuration

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
		return config, err
	}
	return config, nil
}

func createClient(config Configuration) *output.OutputClient {
	return output.NewClient(config.Output.URL, config.Output.Username, config.Output.Password, config.Output.Database, config.Output.Measurement, config.Output.Tags)
}

func recordClimate(dht *dht.DHT, client *output.OutputClient) {
	samples := 3
	humidity := 0.0
	temperature := 0.0
	for i := 0; i < samples; i++ {
		_humidity, _temperature, err := dht.ReadRetry(15)
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		humidity += _humidity
		temperature += _temperature
	}
	humidity /= float64(samples)
	temperature /= float64(samples)
	client.WritePoints(humidity, temperature)
}

func main() {
	config, err := readConfig()
	if err != nil {
		fmt.Println("Configuration Error:", err)
		return
	}

	client := createClient(config)
	defer client.Close()

	err = dht.HostInit()
	if err != nil {
		fmt.Println("HostInit error:", err)
		return
	}

	dht, err := dht.NewDHT("GPIO2", dht.Celsius, "")
	if err != nil {
		fmt.Println("NewDHT error:", err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("Stop recording!\nReceived Signal:", sig)
			os.Exit(1)
		}
	}()

	fmt.Println("Start recording...")
	start := time.Now().UnixNano()
	for {
		recordClimate(dht, client)
		duration := float64(time.Now().UnixNano()-start) / 1000000000.0
		sleeping := config.Input.RefreshInterval - duration
		if sleeping > 0.1 {
			time.Sleep(time.Duration(sleeping*1000) * time.Millisecond)
		}
		start = time.Now().UnixNano()
	}
}
