package main

import (
	"log"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/Pod-Box/swap2p-tg/config"
	"github.com/Pod-Box/swap2p-tg/pkg/bot"
	"github.com/Pod-Box/swap2p-tg/pkg/processor"
	"github.com/Pod-Box/swap2p-tg/pkg/swap2p"
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
	apiClient, err := api.NewClientWithResponses(cfg.Swap2p.GetHost())
	if err != nil {
		log.Fatal(err)
	}
	swapAPI := swap2p.NewClient(&cfg.Swap2p, logger, apiClient)
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
