package go_chipay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/listenfengyang/go-chipay/utils"
)

func (cli *Client) WithdrawReq(req ChipPayWithdrawReq) (*ChipPayWithdrawRsp, error) {
	rawURL := cli.Params.WithdrawURL
	payload, _, err := cli.signWithdraw(req)
	if err != nil {
		return nil, err
	}

	var result ChipPayWithdrawRsp
	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(payload).
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		Post(rawURL)

	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("PSPResty#chipay#withdraw->%s", string(restLog))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		//反序列化错误会在此捕捉
		return nil, fmt.Errorf("status code: %d", resp.StatusCode())
	}

	if resp.Error() != nil {
		//反序列化错误会在此捕捉
		return nil, fmt.Errorf("%v, body:%s", resp.Error(), resp.Body())
	}

	return &result, nil
}

func (cli *Client) signWithdraw(req ChipPayWithdrawReq) (map[string]interface{}, string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, "", err
	}
	var payload map[string]interface{}
	if err = json.Unmarshal(body, &payload); err != nil {
		return nil, "", err
	}

	if req.CompanyID == 0 && cli.Params.CompanyID != 0 {
		payload["companyId"] = cli.Params.CompanyID
	}
	if req.Kyc == "" {
		payload["kyc"] = "2"
	}
	if req.OrderType == 0 {
		payload["orderType"] = 2
	}
	if req.AreaCode == "" {
		payload["areaCode"] = "86"
	}
	if req.CoinSign == "" {
		payload["coinSign"] = "5" //数字货币标识：5.usdt
	}

	// 当payCoinSign为cny时买单支持2.支付宝 , 3.银行卡方式，卖单支持 3.Bank card方式
	// payCoinSign为vnd时买单支持 1.MOMO , 3.Bank card 方式，卖单支持 3.Bank card 方式
	if req.OrderPayChannel == nil {
		payload["orderPayChannel"] = 3
	}
	if req.OrderTime == "" {
		payload["orderTime"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
	if req.AsyncURL == "" && cli.Params.WithdrawAsyncUrl != "" {
		payload["asyncUrl"] = cli.Params.WithdrawAsyncUrl
	}
	if req.SyncURL == "" && cli.Params.SyncURL != "" {
		payload["syncUrl"] = cli.Params.SyncURL
	}

	for key, value := range req.Extra {
		payload[key] = value
	}
	delete(payload, "extra")
	delete(payload, "sign")

	sign, raw, err := utils.SignMap(payload, cli.Params.PrivateKey, "sign")
	if err != nil {
		return nil, "", err
	}
	payload["sign"] = sign
	return payload, raw, nil
}
