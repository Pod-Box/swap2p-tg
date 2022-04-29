package replies

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/enescakir/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/Pod-Box/swap2p-tg/pkg/types"
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
		text = fmt.Sprintf("Oops, something went wrong ☹ \nError: %+v", info[0])
	}
	return &ReplyData{
		text: text,
	}
}

func GetCreatedReplyData(trade *api.Trade, host string, contractType types.ButtonData) *ReplyData {
	urlAPI, _ := url.Parse(host)
	q := urlAPI.Query()
	q.Add("XAssetAddress", trade.XAsset)
	q.Add("YAssetAddress", trade.YAsset)
	q.Add("XAmount", trade.XAmount)
	q.Add("YAmount", trade.YAmount)
	q.Add("contract", string(contractType))
	if trade.YAddress != "" {
		q.Add("YOwner", trade.YAddress)
	}

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

func GetAssetAddedReplyData() *ReplyData {
	return &ReplyData{
		text:   "Asset was added! Now this token will be parsed for all balances" + emoji.MoneyBag.String(),
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
			tgbotapi.NewKeyboardButton(string(types.CreateTradeButton)),
			tgbotapi.NewKeyboardButton(string(types.BrowseTradesButton)),
			tgbotapi.NewKeyboardButton(string(types.AddTokenButton)),
			tgbotapi.NewKeyboardButton(string(types.MyAccountButton)),
		))
}

func GetUserInfoReplyData(data *api.PersonalData) *ReplyData {
	formatted := data.WalletAddress[0:4] + "..." + data.WalletAddress[len(data.WalletAddress)-4:]

	bstring := ""
	for _, balance := range data.Balance {
		bstring += fmt.Sprintf("%v: %v\n", balance.AssetFullName, formatBigInt(balance.Amount, int64(balance.Decimals)))
	}
	return &ReplyData{
		text: fmt.Sprintf("<strong>Your account:</strong>\n<strong>Wallet:</strong> %v\n<strong>Tokens:</strong>\n%v", formatted, bstring),
	}
}

func GetTradeReplyData(trade api.Trade, host string) *ReplyData {
	xAm, _ := decimal.NewFromString(string(trade.XAmount))
	xAm = xAm.Shift(int32(-trade.XDecimals))
	yAm, _ := decimal.NewFromString(string(trade.YAmount))
	yAm = yAm.Shift(int32(-trade.YDecimals))
	shortedXAsset := formatCheckVerificatedAsset(trade.XAsset)
	shortedYAsset := formatCheckVerificatedAsset(trade.YAsset)
	urlAPI, _ := url.Parse(host)
	urlAPI.Path += "/" + strconv.Itoa(trade.Id)
	q := urlAPI.Query()
	q.Add("escrowType", string(trade.Type))
	urlAPI.RawQuery = q.Encode()
	header := "Type: "

	switch trade.Type {
	case api.TradeTypeN2020:
		header += "token" + emoji.RightArrow.String() + "token"
	case api.TradeTypeN20721:
		header += "token" + emoji.RightArrow.String() + "NFT"
	case api.TradeTypeN72120:
		header += "NFT" + emoji.RightArrow.String() + "token"
	case api.TradeTypeN721721:
		header += "NFT" + emoji.RightArrow.String() + "NFT"
	}

	return &ReplyData{
		text: fmt.Sprintf("<strong>%v</strong>\n<strong>%v%v%v</strong>\n%v\n<strong>%v%v%v</strong>\n",
			header, xAm, emoji.Coin, shortedXAsset, emoji.DownArrow, yAm, emoji.Coin, shortedYAsset),
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Accept trade", urlAPI.String()),
			),
		),
	}
}

func GetCreateTradeReplyData() *ReplyData {
	return &ReplyData{
		text: "What type of trade you want to create? " + emoji.ThinkingFace.String(),
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("ERC-20%vERC-20", emoji.RightArrow), string(types.T2T),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("ERC-20%vERC-721", emoji.RightArrow), string(types.T2NFT),
				),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("ERC-721%vERC-20", emoji.RightArrow), string(types.NFT2T),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("ERC-721%vERC-721", emoji.RightArrow), string(types.NFT2NFT),
				),
			),
		),
	}
}

func GetCreateTradeT2TReplyData(balance api.Balance) *ReplyData {
	assetButtons := []tgbotapi.InlineKeyboardButton{}
	for _, asset := range balance {
		assetButtons = append(assetButtons, tgbotapi.NewInlineKeyboardButtonData(string(asset.AssetFullName), string(asset.AssetFullName)))
	}

	return &ReplyData{
		text: "Select asset from the list below",
		markup: tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				assetButtons...,
			)),
	}
}

func GetCreateTradeNFTReplyData() *ReplyData {
	return &ReplyData{
		text: fmt.Sprintf("Provide your NFT address %v", emoji.FramedPicture),
	}
}

func GetCreateTradeNFTYReplyData() *ReplyData {
	return &ReplyData{
		text: fmt.Sprintf("Provide desired NFT address %v", emoji.FramedPicture),
	}
}

func GetAddTokenReplyData() *ReplyData {
	return &ReplyData{
		text: fmt.Sprintf("Provide token address that you want us to know about %v", emoji.LightBulb),
	}
}

func GetIsPersonalReplyData() *ReplyData {
	return &ReplyData{
		text: fmt.Sprintf("If you want your trade to be personal - "+
			"provide target wallet. Otherwise send any text. %v", emoji.LightBulb),
	}
}

func GetCreateTradeAmountReplyData(tradeType types.ButtonData) *ReplyData {
	text := "Now type exact amount you want to trade"
	if tradeType == types.NFT2NFT || tradeType == types.NFT2T {
		text = "Now provide your NFT index"
	}
	return &ReplyData{
		text:   text,
		markup: tgbotapi.NewRemoveKeyboard(false),
	}
}

func GetCreateTradeYAssetReplyData(tradeType types.ButtonData) *ReplyData {
	text := "You're almost done! Provide desired token address"
	if tradeType == types.NFT2NFT || tradeType == types.T2NFT {
		text = "You're almost done! Provide desired NFT address"
	}
	return &ReplyData{
		text: text,
	}
}

func GetCreateTradeYAmountReplyData(tradeType types.ButtonData) *ReplyData {
	text := "Last step :) Type desired token amount"
	if tradeType == types.NFT2NFT || tradeType == types.T2NFT {
		text = "Last step :) Type desired NFT index"
	}
	return &ReplyData{
		text: text,
	}
}

func GetCreateTradeFinishedReplyData(trade *api.Trade, tradeType types.ButtonData) *ReplyData {
	offasset := "You'll trade your token asset:"
	offamount := "Offered token amount:"
	forasset := "For asset token address:"
	foramount := "Desired token amount:"
	pers := "Personal: false"
	switch tradeType {
	case types.NFT2NFT:
		offasset = "You'll trade your NFT:"
		offamount = "Offered NFT index:"
		forasset = "For NFT:"
		foramount = "Desired NFT index:"
	case types.NFT2T:
		offasset = "You'll trade your NFT:"
		offamount = "Offered NFT index:"

	case types.T2NFT:
		forasset = "For NFT:"
		foramount = "Desired NFT index:"
	}

	if trade.YAddress != "" {
		pers = "Personal: " + trade.YAddress
	}

	log.Default().Printf("%+v", trade)

	return &ReplyData{
		text: fmt.Sprintf("Well done! You've created your trade offer!\n"+
			"<strong>%v</strong> %v\n"+
			"<strong>%v</strong> %v\n"+
			"<strong>%v</strong> %v\n"+
			"<strong>%v</strong> %v\n"+
			"<strong>%v</strong>\n"+
			"Do you want to proceed?",
			offasset, formatCheckVerificatedAsset(trade.XAsset),
			forasset, formatCheckVerificatedAsset(trade.YAsset),
			offamount, trade.XAmount,
			foramount, trade.YAmount,
			pers,
		),
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

func formatCheckVerificatedAsset(asset string) string {
	formatted := asset[0:4] + "..." + asset[len(asset)-4:]
	if v, ok := verifiedMap[asset]; ok {
		formatted = fmt.Sprintf("%v%v", v, emoji.CheckMarkButton)
	} else {
		formatted = fmt.Sprintf("%v%v", formatted, emoji.ThinkingFace)
	}
	return formatted
}

var verifiedMap = map[string]string{
	"0x5a87f76aB89916aC92056E646cA93c25bbbb6D88": "TokenX",
	"0x82e2379179Ba2583B8D2d21FdaDd852Ca8Fa1Be1": "TokenY",
	"0xadD10A46e330c0e261e4cC796D7491BCAff632Cb": "SPP",
}
