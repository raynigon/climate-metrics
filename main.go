package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/raynigon/climate-metrics/dht"
	"github.com/raynigon/climate-metrics/output"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Input struct {
		Gpio string `yaml:"gpio"`
	} `yaml:"input"`
	Output struct {
		Url         string            `yaml:"url"`
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
	return output.NewClient(config.Output.Url, config.Output.Username, config.Output.Password, config.Output.Database, config.Output.Measurement, config.Output.Tags)
}

func recordClimate(dht *dht.DHT, client *output.OutputClient) {
	humidity, temperature, err := dht.ReadRetry(15)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
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

	fmt.Println("Start recording...")
	start := time.Now().UnixNano()
	for {
		recordClimate(dht, client)
		duration := float64(time.Now().UnixNano() - start)
		sleeping := 15.0 - duration/1000000000.0
		if sleeping > 0.1 {
			time.Sleep(time.Duration(sleeping*1000) * time.Millisecond)
		}
		start = time.Now().UnixNano()
	}
	fmt.Println("Stop recording!")
}
