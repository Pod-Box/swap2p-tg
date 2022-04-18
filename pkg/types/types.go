package types

type ChatID int64

type Command string

const (
	Start       Command = "start"
	Help        Command = "help"
	Settings    Command = "settings"
	SetAddress  Command = "set_address"
	AcceptOffer Command = "accept_offer"
)

type UserState struct {
	State State
	Step  Step
}

type State string

const (
	StateNew        State = "new"
	StateCreateSwap State = "create_swap"
	StateDefault    State = "default"
	StateAcceptSwap State = "accept_swap"
)

type Step string

const (
	CreateSwapInitStep  Step = "create_swap_init_step"
	CreateSwapOfferStep Step = "create_swap_offer_step"
	CreateSwapLockStep  Step = "create_swap_lock_step"

	AcceptSwapInitStep   Step = "accept_swap_init_step"
	AcceptSwapAmountStep Step = "accept_swap_amount_step"
	AcceptSwapLockStep   Step = "accept_swap_lock_step"

	SecondStep Step = "second_step"
)

type ButtonText string

const (
	CreateTradeButton  ButtonText = "Create new trade"
	BrowseTradesButton ButtonText = "Browse trades"
	MyAccountButton    ButtonText = "My account"
)
