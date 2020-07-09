# Climate Metrics

Records the current temperature and humidity (recorded by a DHT22 / AM2302 / DHT11 sensor)

## Configuration

The Configuration is stored in a `config.yml` which has to be in the working directory of the application

```yaml
input:
    gpio: "GPIO2"
    refreshInterval: 15 # in seconds
output:
    url: "http://my-influx-db:8086"
    username: "my_influx_user"
    password: "my_influx_password"
    database: "my_influx_database"
    measurement: "measurement_name"
    tags:
        host: "tags which will be appended to each Point"
```

## Credits
Sensor Code copied from: https://github.com/MichaelS11/go-dht