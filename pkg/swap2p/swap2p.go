package swap2p

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/IMB-a/swap2p-tg/config"
	"github.com/IMB-a/swap2p-tg/pkg/types"
)

type Client struct {
	cfg *config.Swap2p
}

func NewClient(cfg *config.Swap2p) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) call(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return data, nil
}
func (c *Client) GetDataByChatID(ctx context.Context, id types.ChatID) (*Data, error) {
	body, err := c.call(ctx, http.MethodGet, c.cfg.GetPath(), http.NoBody)
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

func (c *Client) SetUserWallet(ctx context.Context, id types.ChatID) error {
	_, err := c.call(ctx, http.MethodPost, c.cfg.GetPath(), http.NoBody)
	if err != nil {
		return err
	}

	return nil
}

type Data struct {
	Wallet string `json:"wallet"`
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
