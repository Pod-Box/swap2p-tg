package processor

import (
	"context"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/IMB-a/swap2p-tg/pkg/replies"
	"github.com/IMB-a/swap2p-tg/pkg/types"
)

type swap2pAPI interface {
	SetUserWallet(ctx context.Context, id types.ChatID) error
}

type Processor struct {
	api swap2pAPI
}

func NewProcessor(api swap2pAPI) *Processor {
	return &Processor{
		api: api,
	}
}

func (p *Processor) Reply(ctx context.Context, inMsg *tgbotapi.Message) tgbotapi.MessageConfig {
	chatID := types.ChatID(inMsg.Chat.ID)
	r := tgbotapi.NewMessage(int64(chatID), "")
	var replyData *replies.ReplyData

	if cmd := inMsg.Command(); cmd != "" {
		switch types.Command(cmd) {
		case types.Start:
			replyData = replies.GetStartCommandReplyData()
		case types.Help:
		case types.Settings:
		case types.SetAddress:
			if p.IsEVMAddress(inMsg.CommandArguments()) {
				if err := p.api.SetUserWallet(ctx, chatID); err != nil {
					return p.ReplyError(inMsg)
				}
				replyData = replies.GetSuccessAddressSetReplyData()
			}
		}
	}
	if replyData != nil {
		r.Text = replyData.GetText()
		r.ReplyMarkup = replyData.GetMarkup()
		r.ParseMode = "HTML"
	}
	return r
}
func (p *Processor) ReplyError(inMsg *tgbotapi.Message) tgbotapi.MessageConfig {
	r := tgbotapi.NewMessage(inMsg.Chat.ID, "")

	replyData := replies.GetErrorReplyData()
	r.Text = replyData.GetText()

	return r
}

func (p *Processor) ReplyNoAddressError(inMsg *tgbotapi.Message) tgbotapi.MessageConfig {
	r := tgbotapi.NewMessage(inMsg.Chat.ID, "")

	replyData := replies.GetErrorReplyData()
	r.Text = replyData.GetText()

	return r
}

func (p *Processor) ProcessRegistration() tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{}
}

func (p *Processor) IsEVMAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
