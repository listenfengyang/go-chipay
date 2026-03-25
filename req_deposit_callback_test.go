package go_chipay

import (
	"testing"
)

func TestCallback(t *testing.T) {
	vLog := VLog{}
	cli := NewClient(vLog, &ChipPayInitParams{
		MerchantInfo: MerchantInfo{
			CompanyID:         COMPANY_ID,
			PrivateKey:        PRIVATE_KEY,
			CallbackPublicKey: CALLBACK_PUBLIC_KEY,
		},
		DepositURL:       DEPOSIT_URL,
		WithdrawURL:      WITHDRAW_URL,
		QueryDepositURL:  QUERY_DEPOSIT_URL,
		QueryWithdrawURL: QUERY_WITHDRAW_URL,
	})

	err := cli.DepositCallback(GenCallbackRequestDemo(), func(ChipPayDepositCallbackReq) error { return nil })
	if err != nil {
		cli.logger.Errorf("Error:%s", err.Error())
		t.Fatal(err)
	}
}

func GenCallbackRequestDemo() ChipPayDepositCallbackReq {
	// 成功
	req := ChipPayDepositCallbackReq{
		Sign:            "IdnEQIOBP10NmCoKpdEAgJZ2Gl6rgbBdScFuC27dthXVSeIyl5qbdGwESzN1tfkNx9OombzpiuyFIx9xOke/oGKLTQfBHwZl5EzSUU9t4UK/pFqwfA6/ls0tTi5reBx8cxuzgjDzwvT61qhjnxEcKPOrHBL+K5JO+iUjWPE17xo=",
		IntentOrderNo:   "4694051025586177",
		Total:           "300.00000000",
		UnitPrice:       "7.00000000",
		CoinSign:        "5",
		TradeOrderTime:  "2026-03-25 15:17:36",
		SuccessAmount:   "42.85710000",
		CoinAmount:      "42.85710000",
		CompanyOrderNum: "37028219790919996X",
		TradeStatus:     "1",
	}

	return req

	// body, err := json.Marshal(req)
	// if err != nil {
	// 	panic(err)
	// }
	// var payload map[string]interface{}
	// if err = json.Unmarshal(body, &payload); err != nil {
	// 	panic(err)
	// }
	// sign, _, err := utils.SignMap(payload, PRIVATE_KEY, "sign")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("sign callback %s\n", sign)
	// req.Sign = sign
	// return req
}
