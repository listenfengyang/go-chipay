package go_chipay

import (
	"github.com/go-resty/resty/v2"
	"github.com/listenfengyang/go-chipay/utils"
)

type Client struct {
	Params *ChipPayInitParams

	ryClient  *resty.Client
	debugMode bool
	logger    utils.Logger
}

func NewClient(logger utils.Logger, params *ChipPayInitParams) *Client {
	return &Client{
		Params:    params,
		ryClient:  resty.New(),
		debugMode: false,
		logger:    logger,
	}
}

func (cli *Client) SetDebugModel(debugModel bool) {
	cli.debugMode = debugModel
}

func (cli *Client) SetMerchantInfo(merchant MerchantInfo) {
	cli.Params.MerchantInfo = merchant
}
