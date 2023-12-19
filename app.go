package mqtt2prom

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func NewApp() (app *cli.App) {
	app = &cli.App{
		Name:   filepath.Base(os.Args[0]),
		Usage:  "subscribe mqtt topics and push them to prometheus metrics",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "config.yaml", Aliases: []string{"c"}, Usage: "specify config file"},
		},
	}
	return
}

func run(cx *cli.Context) error {
	configFile := cx.String("config")
	cfg, err := NewConfig(configFile)
	if err != nil {
		return err
	}
	if err = StartMQTT(cx.Context, cfg); err != nil {
		return err
	}
	if err = StartProm(cfg); err != nil {
		return err
	}
	select {}
}
