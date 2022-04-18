package types

type ChatID int64

type Command string

const (
	Start      Command = "start"
	Help       Command = "help"
	Settings   Command = "settings"
	SetAddress Command = "set_address"
)
