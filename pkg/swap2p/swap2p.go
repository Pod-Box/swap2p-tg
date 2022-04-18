package swap2p

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/IMB-a/swap2p-tg/config"
	"github.com/IMB-a/swap2p-tg/pkg/types"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	cfg    *config.Swap2p
}

func NewClient(cfg *config.Swap2p, logger *zap.Logger) *Client {
	return &Client{
		cfg:    cfg,
		logger: logger,
	}
}

var udata *Data = &Data{
	Wallet: "",
	State:  "new",
	Step:   "",
	Balance: []Balance{
		{
			Asset:  "",
			Amount: "",
		},
	},
}

func (c *Client) call(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	// req, err := http.NewRequestWithContext(ctx, method, url, body)
	// if err != nil {
	// 	return nil, err
	// }

	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("failed with %v", resp.StatusCode)
	// }

	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// defer resp.Body.Close()
	switch url {
	case "get_data":
		b, _ := json.Marshal(udata)
		return b, nil
	case "set_wallet":
		udata.Wallet = "asdasdasdsaa"
		return []byte{}, nil
	case "user_state":
		return []byte{}, nil
	case "all_trades":
		tr := &TradesList{
			Trades: []*Trade{
				{
					ID:          0,
					OfferAsset:  "BTC",
					WantAsset:   "USDT",
					OfferAmount: "0.01",
					WantAmount:  "15.51",
					Expires:     "01.01.2001",
				},
				{
					ID:          0,
					OfferAsset:  "WAVES",
					WantAsset:   "WX",
					OfferAmount: "100",
					WantAmount:  "15342.69",
					Expires:     "01.02.2002",
				},
			},
		}
		b, _ := json.Marshal(&tr)
		return b, nil
	}
	return []byte{}, nil
}
func (c *Client) GetDataByChatID(ctx context.Context, id types.ChatID) (*Data, error) {
	c.logger.Sugar().Warn(c.cfg.GetDataByChatIDPath())
	body, err := c.call(ctx, http.MethodGet, c.cfg.GetDataByChatIDPath(), http.NoBody)
	if err != nil {
		return nil, err
	}

	data := &Data{}
	if err := data.ParseFromBytes(body); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) IsUserWalletPresents(ctx context.Context, id types.ChatID) (bool, error) {
	data, err := c.GetDataByChatID(ctx, id)
	if err != nil {
		return false, err
	}
	if data != nil {
		if data.GetWallet() != "" {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) SetUserWallet(ctx context.Context, id types.ChatID, wallet string) error {
	_, err := c.call(ctx, http.MethodPost, c.cfg.GetSetWalletPath(), http.NoBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetUserState(ctx context.Context, userState types.UserState) error {
	_, err := c.call(ctx, http.MethodPost, c.cfg.GetSetUserStatePath(), http.NoBody)
	if err != nil {
		return err
	}

	udata.State = string(userState.State)
	udata.Step = string(userState.Step)

	return nil
}

func (c *Client) GetAllTrades(ctx context.Context) (*TradesList, error) {
	c.logger.Info(c.cfg.GetAllTradesPath())
	body, err := c.call(ctx, http.MethodPost, c.cfg.GetAllTradesPath(), http.NoBody)
	if err != nil {
		return nil, err
	}

	trades := &TradesList{}
	if err := trades.ParseFromBytes(body); err != nil {
		return nil, err
	}
	c.logger.Sugar().Warnf("%+v", trades)
	return trades, nil
}

type Data struct {
	Wallet  string    `json:"wallet"`
	State   string    `json:"state"`
	Step    string    `json:"step"`
	Balance []Balance `json:"balances"`
}
type Balance struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

func (d *Data) ParseFromBytes(b []byte) error {
	if err := json.Unmarshal(b, d); err != nil {
		return err
	}
	return nil
}

func (d *Data) GetWallet() string {
	return d.Wallet
}

func (d *Data) GetSessionState() types.UserState {
	return types.UserState{
		State: types.State(d.State),
		Step:  types.Step(d.Step),
	}
}

type TradesList struct {
	Trades []*Trade `json:"trades"`
}

type Trade struct {
	ID          int    `json:"ID"`
	OfferAsset  string `json:"ticker"`
	WantAsset   string `json:"want_asset"`
	OfferAmount string `json:"amount_got"`
	WantAmount  string `json:"amount_want"`
	Expires     string `json:"expires"`
}

func (t *TradesList) ParseFromBytes(b []byte) error {
	if err := json.Unmarshal(b, t); err != nil {
		return err
	}
	return nil
}
