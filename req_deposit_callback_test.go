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
	//{"sign":"EiKv9Zw1KwVaNcRtVqvmaBVbnDfR4HWpdlBySC2DLnb7qFQQOq2SbELbmGORaQqaOxR3ANZbMHukpWEFhSP9QNkmJ2OxJMnBCAFcQrRQ8m7FYqNNVi+5omG4NuRwgDCn+wgw4bwAeEu9UZB0I4CQD1PJxMzU0KuCmNZFZJ2ftfc=",
	// "intentOrderNo":"4694127072989185","total":"38.00000000",
	// "unitPrice":"7.00000000","coinSign":"5","tradeOrderTime":"2026-03-26 11:53:25",
	// "successAmount":"5.42850000","coinAmount":"5.42850000","companyOrderNum":"202603260652380412","tradeStatus":"1"}
	// req := ChipPayDepositCallbackReq{
	// 	Sign:            "EiKv9Zw1KwVaNcRtVqvmaBVbnDfR4HWpdlBySC2DLnb7qFQQOq2SbELbmGORaQqaOxR3ANZbMHukpWEFhSP9QNkmJ2OxJMnBCAFcQrRQ8m7FYqNNVi+5omG4NuRwgDCn+wgw4bwAeEu9UZB0I4CQD1PJxMzU0KuCmNZFZJ2ftfc=",
	// 	IntentOrderNo:   "4694127072989185",
	// 	Total:           "38.00000000",
	// 	UnitPrice:       "7.00000000",
	// 	CoinSign:        "5",
	// 	TradeOrderTime:  "2026-03-26 11:53:25",
	// 	SuccessAmount:   "5.42850000",
	// 	CoinAmount:      "5.42850000",
	// 	CompanyOrderNum: "202603260652380412",
	// 	TradeStatus:     "1",
	// }

	req := ChipPayDepositCallbackReq{
		Sign:            "EiKv9Zw1KwVaNcRtVqvmaBVbnDfR4HWpdlBySC2DLnb7qFQQOq2SbELbmGORaQqaOxR3ANZbMHukpWEFhSP9QNkmJ2OxJMnBCAFcQrRQ8m7FYqNNVi+5omG4NuRwgDCn+wgw4bwAeEu9UZB0I4CQD1PJxMzU0KuCmNZFZJ2ftfc=",
		IntentOrderNo:   "4694127072989185",
		Total:           "38.00000000",
		UnitPrice:       "7.00000000",
		CoinSign:        "5",
		TradeOrderTime:  "2026-03-26 11:53:25",
		SuccessAmount:   "5.42850000",
		CoinAmount:      "5.42850000",
		CompanyOrderNum: "202603260652380412",
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
