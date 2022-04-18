package main

import (
	"log"

	"github.com/IMB-a/swap2p-tg/config"
	"github.com/IMB-a/swap2p-tg/pkg/bot"
	"github.com/IMB-a/swap2p-tg/pkg/processor"
	"github.com/IMB-a/swap2p-tg/pkg/swap2p"
	"go.uber.org/zap"
)

const defaultCfgPath = "./config/"

func main() {
	cfg, err := config.ReadConfig(defaultCfgPath)
	logger, _ := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	logger.Sugar().Infof("%+v", cfg)
	swapAPI := swap2p.NewClient(&cfg.Swap2p, logger)
	proc := processor.NewProcessor(swapAPI)

	tgbot, err := bot.NewBot(cfg.Token, swapAPI, proc, logger)
	if err != nil {
		log.Fatal(err)
	}

	err = tgbot.ListenForUpdates()
	if err != nil {
		log.Fatal(err)
	}
}
