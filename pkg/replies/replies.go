package replies

import (
	"fmt"
	"strconv"

	"github.com/IMB-a/swap2p-tg/pkg/swap2p"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ReplyData struct {
	text   string
	markup interface{}
}

func (r *ReplyData) GetText() string {
	return r.text
}

func (r *ReplyData) GetMarkup() interface{} {
	return r.markup
}

func GetStartCommandReplyData() *ReplyData {
	return &ReplyData{
		text: "✌️ Hello! Tap this button to open Metamask🦊, then paste your wallet address here.\n" +
			"<strong>❗NOTICE:</strong>\n" +
			"This address will be used as receiver in all your trades, so fill it carefully!\n" +
			"You can always change your address via <em>/set_address</em> command.\n",
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Go to 🦊", "https://metamask.app.link"),
			),
		),
	}
}

func GetErrorReplyData() *ReplyData {
	return &ReplyData{
		text: "Oops, something went wrong ☹. Please, try again.",
	}
}

func GetAddressErrorReplyData() *ReplyData {
	return &ReplyData{
		text: "Your address wasn't correctly set. Please, call <em>/set_address</em> or <em>/start</em>.",
	}
}

func GetSuccessAddressSetReplyData() *ReplyData {
	return &ReplyData{
		text:   "Well done, your address was set!",
		markup: GetDefaultButtons(),
	}
}

func GetDefaultButtons() interface{} {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Create new trade"),
			tgbotapi.NewKeyboardButton("Browse trades"),
			tgbotapi.NewKeyboardButton("My account"),
		))
}

func GetUserInfoReplyData(data *swap2p.Data) *ReplyData {
	return &ReplyData{
		text: fmt.Sprintf("<strong>Your account:</strong>\nYour wallet: %v", data.GetWallet()),
	}
}

func GetTradeReplyData(trade *swap2p.Trade) *ReplyData {
	return &ReplyData{
		text: fmt.Sprintf("Offer: <strong>%v %v</strong> for <strong>%v %v</strong>, expires: %v\n",
			trade.OfferAsset, trade.OfferAmount, trade.WantAsset, trade.WantAmount, trade.Expires),
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Accept trade", "accept-trade-"+strconv.Itoa(trade.ID)),
			),
		),
	}
}
