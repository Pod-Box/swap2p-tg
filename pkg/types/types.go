package types

import "fmt"

type ChatID int64

type Command string

const (
	Start       Command = "start"
	Cancel      Command = "cancel"
	Help        Command = "help"
	Settings    Command = "settings"
	SetAddress  Command = "set_address"
	AcceptOffer Command = "accept_offer"
)

type State string

const (
	StateNew       State = "new"
	StateNewSecond State = "new_second"
	StateAddition  State = "addition"
	StateAddToken  State = "add_token"

	StateCreateTrade         State = "create_trade"
	StateCreateTradeType     State = "create_trade_type"
	StateCreateTradeAmount   State = "create_trade_amount"
	StateCreateTradeYAsset   State = "create_trade_y_asset"
	StateCreateTradeYAmount  State = "create_trade_y_amount"
	StateCreateTradeFinished State = "create_trade_finished"
	StateCreateTradePersonal State = "create_trade_personal"
	StateDefault             State = "default"
	StateAcceptTrade         State = "accept_trade"
)

type ButtonText string

const (
	CreateTradeButton           ButtonText = "Create new trade"
	BrowseTradesButton          ButtonText = "Browse trades"
	MyAccountButton             ButtonText = "My account"
	AddTokenButton              ButtonText = "Add token"
	AcceptCreateButton          ButtonText = "Yes, let's go"
	DeclineCreateButton         ButtonText = "No, cancel my order"
	SignCreateTransactionButton ButtonText = "Sign transaction!"
)

type ButtonData string

const (
	T2T     ButtonData = "20_20"
	T2NFT   ButtonData = "20_721"
	NFT2T   ButtonData = "721_20"
	NFT2NFT ButtonData = "721_721"
)

var ErrNotFound = fmt.Errorf("data wasn't found")
var ErrOther = fmt.Errorf("data wasn't found")
