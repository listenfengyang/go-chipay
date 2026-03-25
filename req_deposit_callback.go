package go_chipay

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/listenfengyang/go-chipay/utils"
)

// VerifyBody 使用回调公钥对原始 JSON Body 完整验签。
func (cli *Client) VerifyBody(body []byte) (bool, string, error) {
	return utils.VerifyBody(body, cli.Params.CallbackPublicKey)
}

// DepositCallback 使用结构体回调参数进行验签，并在验签通过后交给业务处理器。
func (cli *Client) DepositCallback(req ChipPayDepositCallbackReq, processor func(ChipPayDepositCallbackReq) error) error {
	// 1) 将结构体回调参数转成 map，确保验签字段和业务字段一一对应。
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

	// 2) sign 字段不参与验签，其他字段必须全部参与。
	delete(payload, "sign")
	ok, raw, err := utils.VerifyMap(payload, sign, cli.Params.CallbackPublicKey)
	if err != nil {
		if cli.logger != nil {
			cli.logger.Errorf("chippay deposit callback verify error, companyOrderNum=%s, intentOrderNo=%s, raw=%s, sign=%s, err=%s", req.CompanyOrderNum, req.IntentOrderNo, raw, sign, err.Error())
		}
		return fmt.Errorf("deposit callback verify failed: %w", err)
	}
	if !ok {
		if cli.logger != nil {
			cli.logger.Errorf("chippay deposit callback verify failed, companyOrderNum=%s, intentOrderNo=%s, raw=%s, sign=%s", req.CompanyOrderNum, req.IntentOrderNo, raw, sign)
		}
		return errors.New("sign verify error")
	}

	// 3) 验签通过后执行业务处理逻辑。
	return processor(req)
}
