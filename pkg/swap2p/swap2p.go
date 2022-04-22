package swap2p

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	api "github.com/Pod-Box/swap2p-backend/api"
	"github.com/Pod-Box/swap2p-tg/config"
	"github.com/Pod-Box/swap2p-tg/pkg/types"
	"go.uber.org/zap"
)

type Client struct {
	c      api.ClientWithResponsesInterface
	logger *zap.Logger
	Cfg    *config.Swap2p
}

func NewClient(cfg *config.Swap2p, logger *zap.Logger, c api.ClientWithResponsesInterface) *Client {
	return &Client{
		c:      c,
		Cfg:    cfg,
		logger: logger,
	}
}

func (c *Client) GetDataByChatID(ctx context.Context, id types.ChatID) (*api.PersonalData, error) {
	ids := strconv.Itoa(int(id))
	resp, err := c.c.GetPersonalDataWithResponse(ctx, api.PChatID(ids))
	if err != nil {
		c.logger.Sugar().Info("GET PERSONAL DATA ERR")
		return nil, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		c.logger.Sugar().Info("NOT FOUND")
		return nil, types.ErrNotFound
	}
	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		c.logger.Sugar().Info("NOT OK OR NIL")
		return nil, types.ErrOther
	}
	return resp.JSON200, nil
}

func (c *Client) InitUserData(ctx context.Context, id types.ChatID) (*api.PersonalData, error) {
	ids := strconv.Itoa(int(id))
	pdata, err := c.GetDataByChatID(ctx, id)
	if errors.Is(err, types.ErrNotFound) {
		c.logger.Sugar().Info(api.PChatID(ids))
		r, err := c.c.InitPersonalDataWithResponse(ctx, api.PChatID(ids))
		if err != nil {
			c.logger.Sugar().Info("INIT PERSONAL DATA ERR", err.Error())
			return nil, err
		}
		if r.StatusCode() != http.StatusOK {
			c.logger.Sugar().Info("NOT OK")
			return nil, fmt.Errorf(string(r.Body))
		}
		pdata, err = c.GetDataByChatID(ctx, id)
		for err == types.ErrNotFound {
			c.logger.Sugar().Info("DIDNT FOUND DATA RETRYING")
			pdata, err = c.GetDataByChatID(ctx, id)
		}
	}

	return pdata, nil
}

func (c *Client) IsUserWalletPresents(ctx context.Context, id types.ChatID) (bool, error) {
	data, err := c.GetDataByChatID(ctx, id)
	if err != nil {
		return false, err
	}
	if data != nil {
		if data.WalletAddress != "" {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) SetUserWallet(ctx context.Context, id types.ChatID, wallet string) error {
	ids := strconv.Itoa(int(id))
	resp, err := c.c.AddWalletWithResponse(
		ctx,
		api.PChatID(ids),
		&api.AddWalletParams{Wallet: api.QWalletAddress(wallet)},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf(string(resp.Body))
	}
	return nil
}

func (c *Client) SetUserState(ctx context.Context, id types.ChatID, userState types.State) error {
	ids := strconv.Itoa(int(id))
	resp, err := c.c.UpdateStateWithResponse(
		ctx,
		api.PChatID(ids),
		&api.UpdateStateParams{State: api.QChatState(userState)},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf(string(resp.Body))
	}
	return nil
}

func (c *Client) GetAllTrades(ctx context.Context) (*api.TradeList, error) {
	offset, limit := api.QOffset(0), api.QLimit(1000)
	resp, err := c.c.GetAllTradesWithResponse(
		ctx,
		&api.GetAllTradesParams{
			Offset: &offset,
			Limit:  &limit,
		},
	)
	if err != nil || resp.StatusCode() != http.StatusOK {
		return nil, err
	}
	trades := resp.JSON200
	c.logger.Sugar().Warnf("%+v", trades)
	return trades, nil
}
