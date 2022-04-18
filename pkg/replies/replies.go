package replies

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
		text: `✌️ Hello! Tap this button to open Metamask 🦊.
        Then paste your wallet address here.
        <strong>❗NOTICE</strong>:
        This address will be used as receiver in all your trades, so fill it carefully!
        You can always change your address via <em>/set_address</em> command.
                `,
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
		text: "Well done, your address was set!",
	}
}
