package go_chipay

import (
	"encoding/json"
	"errors"
)

func (cli *Client) WithdrawCallback(body []byte, processor func(ChipPayWithdrawCallbackReq) error) error {
	ok, _, err := cli.VerifyBody(body)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("sign verify error")
	}

	var req ChipPayWithdrawCallbackReq
	if err = json.Unmarshal(body, &req); err != nil {
		return err
	}
	return processor(req)
}
