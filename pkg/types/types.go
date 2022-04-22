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
	StateNew                 State = "new"
	StateNewSecond           State = "new_second"
	StateAddition            State = "addition"
	StateCreateTrade         State = "create_trade"
	StateCreateTradeAmount   State = "create_trade_amount"
	StateCreateTradeYAsset   State = "create_trade_y_asset"
	StateCreateTradeYAmount  State = "create_trade_y_amount"
	StateCreateTradeFinished State = "create_trade_finished"
	StateDefault             State = "default"
	StateAcceptTrade         State = "accept_trade"
)

type ButtonText string

const (
	CreateTradeButton           ButtonText = "Create new trade"
	BrowseTradesButton          ButtonText = "Browse trades"
	MyAccountButton             ButtonText = "My account"
	MyTradesButton              ButtonText = "My trades"
	AcceptCreateButton          ButtonText = "Yes, let's go"
	DeclineCreateButton         ButtonText = "No, cancel my order"
	SignCreateTransactionButton ButtonText = "Sign transaction!"
)

var ErrNotFound = fmt.Errorf("data wasn't found")
var ErrOther = fmt.Errorf("data wasn't found")
