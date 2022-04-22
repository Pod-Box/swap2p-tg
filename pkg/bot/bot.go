package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/Pod-Box/swap2p-tg/pkg/processor"
	"github.com/Pod-Box/swap2p-tg/pkg/swap2p"
	"github.com/Pod-Box/swap2p-tg/pkg/types"
)

const BotUpdateTimeout = 60

type Handler interface {
	HandleUpdate(tgbotapi.Update) (tgbotapi.MessageConfig, error)
}

type Bot struct {
	token   string
	swapAPI *swap2p.Client
	botAPI  *tgbotapi.BotAPI
	pr      *processor.Processor
	logger  *zap.Logger
}

func NewBot(token string, swapAPI *swap2p.Client, pr *processor.Processor, logger *zap.Logger) (*Bot, error) {
	botapi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botapi.Debug = true

	return &Bot{
		token:   token,
		botAPI:  botapi,
		swapAPI: swapAPI,
		pr:      pr,
		logger:  logger,
	}, nil
}

func (b *Bot) ListenForUpdates() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = BotUpdateTimeout
	ctx := context.Background()
	updates := b.botAPI.GetUpdatesChan(u)
	for update := range updates {
		if err := b.handleUpdate(ctx, &update); err != nil {
			b.logger.Sugar().Error(err)
		}
	}
	return nil
}

func (b *Bot) handleUpdate(ctx context.Context, update *tgbotapi.Update) error {
	var msg []tgbotapi.MessageConfig

	chat := update.FromChat()
	if chat == nil {
		return fmt.Errorf("wtf")
	}
	chatID := types.ChatID(chat.ID)
	data, err := b.swapAPI.InitUserData(ctx, chatID)
	if err != nil {
		b.logger.Sugar().Error(err.Error())
		msg = b.pr.ReplyError(update.Message)
		_, err = b.botAPI.Send(msg[0])
		return err
	}
	if data == nil {
		b.logger.Sugar().Error("no data was found")
		msg = b.pr.ReplyError(update.Message)
		_, err = b.botAPI.Send(msg[0])
		return err
	}
	b.logger.Sugar().Infof("USER_DATA:%+v", data)

	switch {
	case update.Message != nil:
		msg, err = b.pr.Reply(ctx, update.Message, data)
		if err != nil {
			b.logger.Sugar().Error(err.Error())
		}
	case update.CallbackQuery != nil:
		msg = b.pr.ReplyQuery(ctx, update.CallbackQuery, data)
	}
	for _, m := range msg {
		_, err = b.botAPI.Send(m)
	}
	return err
}
