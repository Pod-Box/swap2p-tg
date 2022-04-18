package processor

import (
	"context"
	"fmt"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/IMB-a/swap2p-tg/pkg/replies"
	"github.com/IMB-a/swap2p-tg/pkg/swap2p"
	"github.com/IMB-a/swap2p-tg/pkg/types"
)

type Processor struct {
	api *swap2p.Client
}

func NewProcessor(api *swap2p.Client) *Processor {
	return &Processor{
		api: api,
	}
}

func (p *Processor) Reply(ctx context.Context, inMsg *tgbotapi.Message, userState types.UserState) (msg []tgbotapi.MessageConfig) {
	chatID := types.ChatID(inMsg.Chat.ID)
	msg = []tgbotapi.MessageConfig{}
	fmt.Println(userState)
	switch userState.State {
	case types.StateNew:
		if userState.Step == "" {
			userState.Step = types.SecondStep
			if err := p.api.SetUserState(ctx, userState); err != nil {
				msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
				return msg
			}
			msg = append(msg, buildMessageFromReply(chatID, replies.GetStartCommandReplyData()))
			return msg
		}
		if !p.IsEVMAddress(inMsg.Text) {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		userState.State = types.StateDefault
		userState.Step = ""
		if err := p.api.SetUserState(ctx, userState); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		if err := p.api.SetUserWallet(ctx, chatID, inMsg.Text); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetSuccessAddressSetReplyData()))
	case types.StateCreateSwap:

	case types.StateAcceptSwap:

	case types.StateDefault:
		msg = p.processButton(ctx, inMsg)
	}

	return msg
}

func (p *Processor) processButton(ctx context.Context, inMsg *tgbotapi.Message) (msg []tgbotapi.MessageConfig) {
	chatID := types.ChatID(inMsg.Chat.ID)
	msg = []tgbotapi.MessageConfig{}

	switch types.ButtonText(inMsg.Text) {
	case types.MyAccountButton:
		data, err := p.api.GetDataByChatID(ctx, chatID)
		if err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetUserInfoReplyData(data)))
	case types.BrowseTradesButton:
		trades, err := p.api.GetAllTrades(ctx)
		if err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		for _, trade := range trades.Trades {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetTradeReplyData(trade)))
		}
	case types.CreateTradeButton:
	}

	return msg
}

func (p *Processor) ReplyQuery(
	ctx context.Context,
	inQuery *tgbotapi.CallbackQuery,
	userState types.UserState,
) (msg []tgbotapi.MessageConfig) {
	//	chatID := types.ChatID(inQuery.From.ID)

	switch userState.State {

	case types.StateDefault:

	case types.StateCreateSwap:

	case types.StateAcceptSwap:

	}

	return
}

func buildMessageFromReply(chatID types.ChatID, r *replies.ReplyData) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(int64(chatID), "")
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = r.GetMarkup()
	msg.Text = r.GetText()
	return msg
}

func (p *Processor) ReplyError(inMsg *tgbotapi.Message) []tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(inMsg.Chat.ID, "")

	replyData := replies.GetErrorReplyData()
	msg.Text = replyData.GetText()

	return []tgbotapi.MessageConfig{msg}
}

func (p *Processor) ReplyNoAddressError(inMsg *tgbotapi.Message) []tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(inMsg.Chat.ID, "")

	replyData := replies.GetAddressErrorReplyData()
	msg.Text = replyData.GetText()

	return []tgbotapi.MessageConfig{msg}
}

func (p *Processor) IsEVMAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
