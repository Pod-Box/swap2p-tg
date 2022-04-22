package replies

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/Pod-Box/swap2p-tg/pkg/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
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

func GetErrorReplyData(info ...string) *ReplyData {
	text := ("Oops, something went wrong ☹. Please, try again")
	if len(info) != 0 {
		text = fmt.Sprintf("Oops, something went wrong ☹. Please, try again.%+v", info[0])
	}
	return &ReplyData{
		text: text,
	}
}

func GetAcceptOfferReplyData(data string, host string) *ReplyData {
	data = strings.TrimPrefix(data, "accept-trade-")
	urlAPI, _ := url.Parse(host)
	q := urlAPI.Query()
	q.Add("escrowIndex", data)
	urlAPI.RawQuery = q.Encode()

	return &ReplyData{
		text: "Click this button to accept this trade!",
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Go to 🦊", urlAPI.String()),
			),
		),
	}
}

func GetCreatedReplyData(trade *api.Trade, host string) *ReplyData {
	urlAPI, _ := url.Parse(host)
	q := urlAPI.Query()
	q.Add("XAssetAddress", trade.XAddress)
	q.Add("YAssetAddress", trade.YAddress)
	q.Add("XAmount", trade.XAmount)
	q.Add("YAmount", trade.YAmount)

	urlAPI.RawQuery = q.Encode()
	urlAPI.Path += "/create"

	return &ReplyData{
		text: "Now click this button to list your offer!",
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(string(types.SignCreateTransactionButton), urlAPI.String()),
			),
		),
	}
}

func GetCancelledReplyData() *ReplyData {
	return &ReplyData{
		text:   "Order was deleted",
		markup: GetDefaultButtons(),
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

func GetDefaultReplyData() *ReplyData {
	return &ReplyData{
		text:   "Returned to the menu >:)",
		markup: GetDefaultButtons(),
	}
}

func GetDefaultButtons() interface{} {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Create new trade"),
			tgbotapi.NewKeyboardButton("Browse trades"),
			tgbotapi.NewKeyboardButton("My trades"),
			tgbotapi.NewKeyboardButton("My account"),
		))
}

func GetUserInfoReplyData(data *api.PersonalData) *ReplyData {
	bstring := ""
	for _, balance := range data.Balance {
		bstring += fmt.Sprintf("Token %v: %v\n", balance.Asset, formatBigInt(balance.Amount, int64(balance.Decimals)))
	}
	return &ReplyData{
		text: fmt.Sprintf("<strong>Your account:</strong>\nWallet: %v\n%v", data.WalletAddress, bstring),
	}
}

func GetTradeReplyData(trade api.Trade) *ReplyData {
	xAm, _ := decimal.NewFromString(string(trade.XAmount))
	yAm, _ := decimal.NewFromString(string(trade.XAmount))
	shortedXAsset := trade.XAsset[0:4] + "..." + trade.XAsset[len(trade.XAsset)-4:len(trade.XAddress)]
	shortedYAsset := trade.YAsset[0:4] + "..." + trade.YAsset[len(trade.YAsset)-4:len(trade.YAddress)]

	return &ReplyData{
		text: fmt.Sprintf("Offer: <strong>%v of %v</strong> for <strong>%v of %v</strong>, expires: %v\n",
			xAm, shortedXAsset, yAm, shortedYAsset, time.Now().Format(time.RFC822)),
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Accept trade", "accept-trade-"+strconv.Itoa(trade.Id)),
			),
		),
	}
}

func GetCreateTradeReplyData(balance api.Balance) *ReplyData {
	assetButtons := []tgbotapi.InlineKeyboardButton{}
	for _, asset := range balance {
		assetButtons = append(assetButtons, tgbotapi.NewInlineKeyboardButtonData(string(asset.Asset), string(asset.Asset)))
	}

	return &ReplyData{
		text: "Select asset from the list below",
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				assetButtons...,
			)),
	}
}

func GetCreateTradeAmountReplyData() *ReplyData {
	return &ReplyData{
		text:   "Now type exact amount you want to trade",
		markup: tgbotapi.NewRemoveKeyboard(false),
	}
}

func GetCreateTradeYAssetReplyData() *ReplyData {
	return &ReplyData{
		text: "You're almost done! Provide desired token address",
	}
}

func GetCreateTradeYAmountReplyData() *ReplyData {
	return &ReplyData{
		text: "Last step :) Type desired token amount",
	}
}

func GetCreateTradeFinishedReplyData(trade *api.Trade) *ReplyData {
	return &ReplyData{
		text: fmt.Sprintf("Well done! You've created your trade offer!\n"+
			"<strong>You'll trade asset:</strong> %v\n"+
			"<strong>For asset address:</strong> %v\n"+
			"<strong>Offered amount:</strong> %v\n"+
			"<strong>Desired amount:</strong> %v\n"+
			"Do you want to proceed?",
			trade.XAsset, trade.YAddress,
			formatBigInt(string(trade.XAmount), int64(trade.XDecimals)), formatBigInt(string(trade.YAmount), 18)),
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(string(types.AcceptCreateButton), string(types.AcceptCreateButton)),
				tgbotapi.NewInlineKeyboardButtonData(string(types.DeclineCreateButton), string(types.DeclineCreateButton)),
			)),
	}
}

func formatBigInt(n string, decimals int64) string {
	amount, _ := decimal.NewFromString(n)
	amount = amount.Shift(-int32(decimals))
	return amount.String()
}
