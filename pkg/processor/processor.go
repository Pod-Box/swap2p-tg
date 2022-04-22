package processor

import (
	"context"
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/Pod-Box/swap2p-tg/pkg/replies"
	"github.com/Pod-Box/swap2p-tg/pkg/swap2p"
	"github.com/Pod-Box/swap2p-tg/pkg/types"
)

type Processor struct {
	api *swap2p.Client
	m   map[types.ChatID]*api.Trade
}

func NewProcessor(swap2pAPI *swap2p.Client) *Processor {
	return &Processor{
		api: swap2pAPI,
		m:   make(map[types.ChatID]*api.Trade),
	}
}

func (p *Processor) Reply(ctx context.Context, inMsg *tgbotapi.Message, data *api.PersonalData) (msg []tgbotapi.MessageConfig, err error) {
	chatID := types.ChatID(inMsg.Chat.ID)
	msg = []tgbotapi.MessageConfig{}
	state := types.State(data.State)

	if inMsg.Command() == string(types.Cancel) {
		if data.WalletAddress == "" {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData("DOLBOEB VVEDI ADDRESS")))
			return msg, err
		}
		if err := p.api.SetUserState(ctx, chatID, types.StateDefault); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetDefaultReplyData()))
		return msg, err
	}

	switch state {
	case types.StateNew:
		state = types.StateNewSecond
		if err := p.api.SetUserState(ctx, chatID, state); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetStartCommandReplyData()))
	case types.StateNewSecond:
		if inMsg.Command() == string(types.Start) {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetStartCommandReplyData()))
			return msg, nil
		}
		if !p.IsEVMAddress(inMsg.Text) {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData("Please provide correct address")))
			return msg, err
		}
		state = types.StateDefault
		if err := p.api.SetUserState(ctx, chatID, state); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}
		if err := p.api.SetUserWallet(ctx, chatID, inMsg.Text); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetSuccessAddressSetReplyData()))
	case types.StateCreateTradeAmount:
		trade, ok := p.m[chatID]
		if !ok {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, fmt.Errorf("no metadata")
		}

		tokenInfo := p.GetTokenFromBalanceByAddress(trade.XAddress, data.Balance)
		if tokenInfo == nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, fmt.Errorf("couldn't find token %v", trade.XAsset)
		}

		gotAmount, err := p.TextToAmount(inMsg.Text, int64(tokenInfo.Decimals))
		if err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}

		if !p.IsSufficientFunds(tokenInfo, gotAmount) {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData("unsufficient funds")))
			return msg, fmt.Errorf("unsufficient funds")
		}

		trade.XAddress = tokenInfo.Address
		trade.XAmount = gotAmount.String()
		trade.XDecimals = tokenInfo.Decimals
		p.m[chatID] = trade

		state = types.StateCreateTradeYAsset
		if err := p.api.SetUserState(ctx, chatID, state); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}

		msg = append(msg, buildMessageFromReply(chatID, replies.GetCreateTradeYAssetReplyData()))
	case types.StateCreateTradeYAsset:
		if !p.IsEVMAddress(inMsg.Text) {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, fmt.Errorf("not an EVM address")
		}

		trade := p.m[chatID]
		trade.YAddress = inMsg.Text
		p.m[chatID] = trade

		state = types.StateCreateTradeYAmount
		if err := p.api.SetUserState(ctx, chatID, state); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetCreateTradeYAmountReplyData()))
	case types.StateCreateTradeYAmount:
		gotAmount, err := p.TextToAmount(inMsg.Text, 18)
		if err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}
		trade := p.m[chatID]
		trade.YAmount = gotAmount.String()
		p.m[chatID] = trade

		state = types.StateCreateTradeFinished
		if err := p.api.SetUserState(ctx, chatID, state); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg, err
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetCreateTradeFinishedReplyData(p.m[chatID])))
	case types.StateDefault:
		msg = p.processButton(ctx, inMsg, data)
	}
	return msg, nil
}

func (p *Processor) processButton(ctx context.Context, inMsg *tgbotapi.Message, data *api.PersonalData) (msg []tgbotapi.MessageConfig) {
	chatID := types.ChatID(inMsg.Chat.ID)
	msg = []tgbotapi.MessageConfig{}

	switch types.ButtonText(inMsg.Text) {
	case types.MyAccountButton:
		msg = append(msg, buildMessageFromReply(chatID, replies.GetUserInfoReplyData(data)))
	case types.BrowseTradesButton:
		trades, err := p.api.GetAllTrades(ctx)
		if err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		if trades == nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData("No trades was found :(")))
			return msg
		}
		for _, trade := range *trades {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetTradeReplyData(trade)))
		}
	case types.CreateTradeButton:
		if len(data.Balance) == 0 {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData("You have no tokens :(")))
			return msg
		}
		if err := p.api.SetUserState(ctx, chatID, types.StateCreateTrade); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		p.m[chatID] = &api.Trade{
			Id:     int(time.Now().Unix()),
			Closed: false,
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetCreateTradeReplyData(data.Balance)))
	}

	return msg
}

func (p *Processor) processCreateTradeButton(ctx context.Context, chatID types.ChatID, data string) (msg []tgbotapi.MessageConfig) {
	msg = []tgbotapi.MessageConfig{}
	trade, ok := p.m[chatID]
	if !ok {
		msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
		return msg
	}
	switch types.ButtonText(data) {
	case types.AcceptCreateButton:
		msg = append(msg, buildMessageFromReply(chatID, replies.GetCreatedReplyData(trade, p.api.Cfg.GetRedirectHost())))
	case types.DeclineCreateButton:
		msg = append(msg, buildMessageFromReply(chatID, replies.GetCancelledReplyData()))
	}
	if err := p.api.SetUserState(ctx, chatID, types.StateDefault); err != nil {
		msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
		return msg
	}
	delete(p.m, chatID)
	return msg
}

func (p *Processor) ReplyQuery(
	ctx context.Context,
	inQuery *tgbotapi.CallbackQuery,
	data *api.PersonalData,
) (msg []tgbotapi.MessageConfig) {
	chatID := types.ChatID(inQuery.From.ID)
	state := types.State(data.State)
	switch state {
	case types.StateCreateTrade:
		tokenInfo := p.GetTokenFromBalance(inQuery.Data, data.Balance)
		if tokenInfo == nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData("No such token!")))
			return msg
		}
		trade := p.m[chatID]
		trade.XAsset = tokenInfo.Asset
		trade.XAddress = tokenInfo.Address
		p.m[chatID] = trade

		state = types.StateCreateTradeAmount
		if err := p.api.SetUserState(ctx, chatID, state); err != nil {
			msg = append(msg, buildMessageFromReply(chatID, replies.GetErrorReplyData()))
			return msg
		}
		msg = append(msg, buildMessageFromReply(chatID, replies.GetCreateTradeAmountReplyData()))
	case types.StateDefault:
		switch inQuery.Data {
		case string(types.ButtonText(types.CreateTradeButton)):
			msg = append(msg, buildMessageFromReply(chatID, replies.GetCreatedReplyData(p.m[chatID], p.api.Cfg.GetRedirectHost())))
		case string(types.ButtonText(types.DeclineCreateButton)):
			msg = append(msg, buildMessageFromReply(chatID, replies.GetCancelledReplyData()))
		default:
			msg = append(msg, buildMessageFromReply(chatID, replies.GetAcceptOfferReplyData(inQuery.Data, p.api.Cfg.GetRedirectHost())))
		}
		return msg
	case types.StateCreateTradeFinished:
		msg = p.processCreateTradeButton(ctx, chatID, inQuery.Data)
	}
	return msg
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

func (p *Processor) IsValidAsset(text string, balances api.Balance) bool {
	for _, balance := range balances {
		if balance.Asset == text {
			return true
		}
	}
	return false
}

func (p *Processor) IsSufficientFunds(tokenInfo *api.SingleBalance, amount *decimal.Decimal) bool {
	am, _ := decimal.NewFromString(string(tokenInfo.Amount))
	return am.GreaterThanOrEqual(*amount)
}

func (p *Processor) TextToAmount(text string, decimals int64) (*decimal.Decimal, error) {
	n, err := decimal.NewFromString(text)
	if err != nil {
		return nil, err
	}
	n = n.Shift(int32(decimals))
	return &n, nil
}

func (p *Processor) GetTokenFromBalance(text string, b api.Balance) *api.SingleBalance {
	for _, v := range b {
		if text == v.Asset {
			return &v
		}
	}
	return nil
}

func (p *Processor) GetTokenFromBalanceByAddress(address string, b api.Balance) *api.SingleBalance {
	for _, v := range b {
		if address == v.Address {
			return &v
		}
	}
	return nil
}
