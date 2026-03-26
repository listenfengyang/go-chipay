package go_chipay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/listenfengyang/go-chipay/utils"
)

func (cli *Client) Deposit(req ChipPayDepositReq) (*ChipPayDepositRsp, error) {
	rawURL := cli.Params.DepositURL
	req.AsyncURL = cli.Params.AsyncURL   // 异步通知地址
	req.CompanyID = cli.Params.CompanyID // 商户ID
	req.SyncURL = cli.Params.SyncURL     // 同步返回地址
	// req.AreaCode = "90"                  //"86"                  // 区号 默认86

	payload, rawSign, err := cli.signStruct(req)
	if err != nil {
		return nil, err
	}

	// payload["companyId"] = cli.Params.CompanyID // 商户ID
	// payload["asyncUrl"] = cli.Params.AsyncURL   // 异步通知地址
	// payload["syncUrl"] = cli.Params.SyncURL     // 同步返回地址
	// payload["areaCode"] = "86"                  // 区号 默认86

	var result ChipPayDepositRsp
	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(payload).
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		Post(rawURL)

	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("PSPResty#chipay#deposit->%s", string(restLog))

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

	_ = rawSign
	return &result, nil
}

func (cli *Client) QueryDeposit(req ChipPayDepositQueryReq) (*ChipPayDepositQueryRsp, error) {
	rawURL := cli.Params.QueryDepositURL
	payload, _, err := cli.signStruct(req)
	if err != nil {
		return nil, err
	}

	var result ChipPayDepositQueryRsp
	resp, err := cli.ryClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetCloseConnection(true).
		R().
		SetBody(payload).
		SetHeaders(getHeaders()).
		SetDebug(cli.debugMode).
		SetResult(&result).
		Post(rawURL)

	restLog, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(utils.GetRestyLog(resp))
	cli.logger.Infof("PSPResty#chipay#deposit->%s", string(restLog))

	if err != nil {
		return nil, err
	}

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

func (cli *Client) signStruct(req interface{}) (map[string]interface{}, string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, "", err
	}
	var payload map[string]interface{}
	if err = json.Unmarshal(body, &payload); err != nil {
		return nil, "", err
	}
	delete(payload, "sign")
	sign, raw, err := utils.SignMap(payload, cli.Params.PrivateKey)
	if err != nil {
		return nil, "", err
	}
	payload["sign"] = sign
	return payload, raw, nil
}

func (cli *Client) logResty(tag string, resp interface{}) {
	if cli.logger == nil {
		return
	}
	r, ok := resp.(interface{ String() string })
	if ok && r != nil {
		cli.logger.Infof("%s->%s", tag, r.String())
	}
}
