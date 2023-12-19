package mqtt2prom

import (
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type metricEntry struct {
	sync.Mutex
	Name   string
	Gauge  prometheus.Gauge
	Active time.Time
}

var metrics sync.Map

func StartProm(cfg *Config) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("gocron.NewScheduler got error", "error", err.Error())
		return err
	}
	if _, err = s.NewJob(gocron.DurationJob(cfg.PushInterval), gocron.NewTask(pushJob, cfg.PushGateway, cfg.PushJob)); err != nil {
		slog.Error("add push job got error", "error", err.Error())
		return err
	}
	if _, err = s.NewJob(gocron.DurationJob(cfg.CleanInterval), gocron.NewTask(cleanJob)); err != nil {
		slog.Error("add clean job got error", "error", err.Error())
		return err
	}
	s.Start()
	return nil
}

func AddOrUpdate(topic string, val float64) {
	m, loaded := metrics.LoadOrStore(topic, &metricEntry{
		Name: topic,
		Gauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: strings.ReplaceAll(topic, "/", "_"),
		}),
		Active: time.Now(),
	})
	if !loaded {
		return
	}
	entry, ok := m.(*metricEntry)
	if !ok {
		metrics.Delete(topic)
		return
	}
	entry.Lock()
	defer entry.Unlock()
	entry.Active = time.Now()
	entry.Gauge.Set(val)
}

func pushJob(url, name string) {
	slog.Debug("prom push job")
	mm := collectMetrics()
	if len(mm) == 0 {
		slog.Info("no metric to push")
		return
	}
	pusher := push.New(url, name)
	for _, cl := range mm {
		pusher.Collector(cl)
	}
	if err := pusher.Push(); err != nil {
		slog.Error("push metircs error", "error", err.Error())
	}
}

func cleanJob(dur time.Duration) {
	slog.Debug("prom clean job")
	now := time.Now()
	metrics.Range(func(key, value any) bool {
		entry, ok := value.(*metricEntry)
		if !ok {
			metrics.Delete(key)
			return true
		}
		entry.Lock()
		defer entry.Unlock()
		if now.Sub(entry.Active) > dur {
			metrics.Delete(key)
		}
		return true
	})
}

func collectMetrics() []prometheus.Collector {
	cc := []prometheus.Collector{}
	metrics.Range(func(key, value any) bool {
		entry, ok := value.(*metricEntry)
		if !ok {
			metrics.Delete(key)
			return true
		}
		cc = append(cc, entry.Gauge)
		return true
	})
	return cc
}
