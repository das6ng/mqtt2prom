package mqtt2prom

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Broker        string        `yaml:"mqtt_broker_url"`
	ClientID      string        `yaml:"mqtt_client_id"`
	Topics        []string      `yaml:"topics"`
	Ignores       []string      `yaml:"ignore_topics"`
	PushGateway   string        `yaml:"pushgateway_url"`
	PushJob       string        `yaml:"push_job_name"`
	PushInterval  time.Duration `yaml:"push_interval"`
	CleanInterval time.Duration `yaml:"clean_interval"`
	CleanDur      time.Duration `yaml:"clean_duration"`
	LogLevel      string        `yaml:"log_level"`

	IgnoredTopics map[string]struct{} `yaml:"-"`
}

func NewConfig(filename string) (c *Config, err error) {
	c = &Config{}
	if bs, err := os.ReadFile(filename); err != nil {
		log.Fatalln("open config file:", err)
	} else if err = yaml.Unmarshal(bs, c); err != nil {
		log.Fatalln("unmarshal config:", err)
	}
	// setup logger
	var lv = slog.LevelError
	if c.LogLevel != "" {
		if err := lv.UnmarshalText([]byte(c.LogLevel)); err != nil {
			log.Println("parse log level err, set to ERROR for default: ", c.LogLevel)
			lv = slog.LevelError
		}
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: lv,
	})))
	// transform ignore topics
	c.IgnoredTopics = lo.SliceToMap(c.Ignores, func(v string) (string, struct{}) {
		return v, struct{}{}
	})
	return
}
