package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/IMB-a/swap2p-tg/pkg/processor"
	"github.com/IMB-a/swap2p-tg/pkg/types"
)

const BotUpdateTimeout = 60

type Handler interface {
	HandleUpdate(tgbotapi.Update) (tgbotapi.MessageConfig, error)
}

type swap2pAPI interface {
	SetUserWallet(ctx context.Context, id types.ChatID) error
	IsUserWalletPresents(ctx context.Context, id types.ChatID) (bool, error)
}

type Bot struct {
	token   string
	swapAPI swap2pAPI
	botAPI  *tgbotapi.BotAPI
	pr      *processor.Processor
	logger  *zap.Logger
}

func NewBot(token string, swapAPI swap2pAPI, pr *processor.Processor, logger *zap.Logger) (*Bot, error) {
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
	var msg tgbotapi.MessageConfig

	chat := update.FromChat()
	chatID := types.ChatID(chat.ID)
	ok, err := b.swapAPI.IsUserWalletPresents(ctx, chatID)
	if err != nil {
		msg = b.pr.ReplyError(update.Message)
	}
	if !ok {
		msg = b.pr.ReplyNoAddressError(update.Message)
	}

	if update.Message != nil {
		msg = b.pr.Reply(ctx, update.Message)
	}

	_, err = b.botAPI.Send(msg)

	return err
}
