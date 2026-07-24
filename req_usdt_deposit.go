package go_chipay

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/listenfengyang/go-chipay/utils"
)

// USDTDeposit 发起 USDT 轮换地址模式入金下单（1.1.3）。
//
// rechargeCoin 传币种编号字符串，USDT 对应 "5"。
// rechargeAmount 为选填，传 nil 表示不限定存币数量。
// 收银台链接在响应 Data.Link 中；各链收款地址在 Data.OrderAddressesDTO.Addresses 中。
//
// 当 rechargeCoin 支持多链时，收银台默认展示所有支持的链类型充值地址。
// 如需在收银台只显示指定链，可设置 req.CoinType，例如 "erc20"、"trc20"、"bep20"、"spl"，
// SDK 将自动在返回的 Link 末尾追加 &cointype=<value> 参数。
func (cli *Client) USDTDeposit(req ChipPayUSDTDepositReq) (*ChipPayUSDTDepositRsp, error) {
	rawURL := cli.Params.USDTDepositURL
	if rawURL == "" {
		return nil, errors.New("USDTDepositURL is not configured")
	}

	// 自动填充商户公共字段
	req.CompanyID = cli.Params.CompanyID
	req.AsyncURL = cli.Params.USDTDepositAsyncUrl
	if req.SyncURL == "" {
		req.SyncURL = cli.Params.SyncURL
	}

	payload, _, err := cli.signStruct(req)
	if err != nil {
		return nil, err
	}

	var result ChipPayUSDTDepositRsp
	resp, err := cli.ryClient.
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(payload).
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		Post(rawURL)

	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("PSPResty#chipay#usdtDeposit->%s", string(restLog))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode())
	}
	if resp.Error() != nil {
		return nil, fmt.Errorf("%v, body:%s", resp.Error(), resp.Body())
	}

	// 若调用方指定了链类型，在收银台 Link 末尾追加 &cointype=<value>，
	// 使收银台仅显示对应链的充值地址。
	if req.CoinType != "" && result.Data.Link != "" {
		result.Data.Link = result.Data.Link + "&cointype=" + req.CoinType
	}

	return &result, nil
}

// USDTDepositCallback 对 USDT 轮换地址模式的回调通知进行验签，验签通过后调用业务处理器。
//
// 验签规则：将回调所有字段（除 sign 外）按字段名升序排列，拼接成
// key=value&key=value 格式字符串，再用平台公钥（SHA256WithRSA）验签。
func (cli *Client) USDTDepositCallback(req ChipPayUSDTDepositCallbackReq, processor func(ChipPayUSDTDepositCallbackReq) error) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	var payload map[string]interface{}
	if err = json.Unmarshal(body, &payload); err != nil {
		return err
	}

	sign := req.Sign
	if sign == "" {
		return errors.New("missing sign field")
	}

	// sign 字段本身不参与验签
	delete(payload, "sign")

	ok, raw, err := utils.VerifyMap(payload, sign, cli.Params.CallbackPublicKey)
	if err != nil {
		if cli.logger != nil {
			cli.logger.Errorf("chippay usdt deposit callback verify error, companyOrderNo=%s, rechargeDetailOrderNo=%s, raw=%s, sign=%s, err=%s",
				req.CompanyOrderNo, req.RechargeDetailOrderNo, raw, sign, err.Error())
		}
		return fmt.Errorf("usdt deposit callback verify failed: %w", err)
	}
	if !ok {
		if cli.logger != nil {
			cli.logger.Errorf("chippay usdt deposit callback verify failed, companyOrderNo=%s, rechargeDetailOrderNo=%s, raw=%s, sign=%s",
				req.CompanyOrderNo, req.RechargeDetailOrderNo, raw, sign)
		}
		return errors.New("sign verify error")
	}

	return processor(req)
}
