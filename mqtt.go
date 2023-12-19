package mqtt2prom

import (
	"context"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/samber/lo"
)

func StartMQTT(ctx context.Context, cfg *Config) error {
	u, err := url.Parse(cfg.Broker)
	if err != nil {
		slog.Error("parse mqtt broker url", "err", err.Error())
		return err
	}
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:     []*url.URL{u},
		KeepAlive:      20,
		OnConnectionUp: func(_ *autopaho.ConnectionManager, _ *paho.Connack) { slog.Info("mqtt connection up") },
		OnConnectError: func(err error) { slog.Error("error whilst attempting connection", "error", err.Error()) },
		ClientConfig: paho.ClientConfig{
			ClientID:      cfg.ClientID,
			Router:        paho.NewSingleHandlerRouter(func(p *paho.Publish) { handleMQTT(cfg, p) }),
			OnClientError: func(err error) { slog.Error("client error", "error", err.Error()) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					slog.Warn("server requested disconnect", "reason", d.Properties.ReasonString)
				} else {
					slog.Warn("server requested disconnect", "reason_code", d.ReasonCode)
				}
			},
		},
	}
	cm, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		slog.Error("establish mqtt connection error", "error", err.Error())
		return err
	}
	if err = cm.AwaitConnection(ctx); err != nil {
		slog.Error("wait mqtt connection error", "error", err.Error())
		return err
	}
	if _, err = cm.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: lo.Map(cfg.Topics, func(t string, _ int) paho.SubscribeOptions {
			return paho.SubscribeOptions{Topic: t, QoS: 1}
		}),
	}); err != nil {
		slog.Error("subscribe error", "error", err.Error())
		return err
	}
	return nil
}

func handleMQTT(cfg *Config, p *paho.Publish) {
	if _, ok := cfg.IgnoredTopics[p.Topic]; ok {
		return
	}
	slog.Debug("new message", "topic", p.Topic, "content", string(p.Payload))
	str := string(p.Payload)
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		slog.Error("parse topic value error", "topic", p.Topic, "value", str, "error", err.Error())
		return
	}
	AddOrUpdate(p.Topic, val)
}
