# mqtt2prom

subscribe mqtt topics and push them to prometheus metrics.

usage:
```
USAGE:
   mqtt2prom [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  specify config file (default: "config.yaml")
   --help, -h                show help
```

config:

```yaml
# mqtt broker url
mqtt_broker_url: mqtt://localhost:1883

# mqtt client id
mqtt_client_id: test_mqtt2prom

# subscribe topics/topic patterns
topics:
  # normal topic
  - "test_topic"
  # topic pattern, I've tested it on mosquitto
  - "/env/#"

# ignore topics, which wont push to prom push gateway
# MUST list one by one, no pattern
ignore_topics:
  - "/env/debug"

# prom pushgateway url
pushgateway_url: http://localhost:9091

# prom push job name
push_job_name: "mqtt2prom"

# prom push job interval
push_interval: 10s

# prom clean job interval
clean_interval: 1m

# clean job will delete metrics which is no data in this duration.
clean_duration: 1h

# log level
log_level: DEBUG
```
