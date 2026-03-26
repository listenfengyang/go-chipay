package go_chipay

import (
	"testing"
)

func TestWithdrawCallback(t *testing.T) {
	vLog := VLog{}
	cli := NewClient(vLog, &ChipPayInitParams{
		MerchantInfo: MerchantInfo{
			CompanyID:         COMPANY_ID,
			PrivateKey:        PRIVATE_KEY,
			CallbackPublicKey: CALLBACK_PUBLIC_KEY,
			WithdrawAsyncUrl:  WITHDRAW_ASYNC_URL,
		},
		DepositURL:       DEPOSIT_URL,
		WithdrawURL:      WITHDRAW_URL,
		QueryDepositURL:  QUERY_DEPOSIT_URL,
		QueryWithdrawURL: QUERY_WITHDRAW_URL,
	})

	err := cli.WithdrawCallback(GenWdCallbackRequestDemo(), func(ChipPayWithdrawCallbackReq) error { return nil })
	if err != nil {
		cli.logger.Errorf("Error:%s", err.Error())
		t.Fatal(err)
	}
}

func GenWdCallbackRequestDemo() ChipPayWithdrawCallbackReq {
	// return ChipPayWithdrawCallbackReq{
	// 	Sign:            "espjodO5isYP8DvcCwJLAYRNInO0VPhwXufYQfcRw2b1UzdJBOaCOpIUFyn4g2dJDuzUcp9jAZqdPrdN/x9ldyNa3zXVG1B3NZsKemxQEt98DrVaD01UqH7+YMENfvr7ATuR6SdwI7pRcHNxWFxJse8LNZmDtzyLzQTi69sCUo4=",
	// 	CompanyOrderNum: "2029884682646034",
	// 	OtcOrderNum:     "4548687456245761_17744988807694",
	// 	CoinAmount:      "38.96100000",
	// 	CoinSign:        "usdt",
	// 	OrderType:       "2",
	// 	TradeStatus:     "1",
	// 	TradeOrderTime:  "2026-03-26 12:21:21",
	// 	UnitPrice:       "7.70000000",
	// 	Total:           "300.00000000",
	// 	SuccessAmount:   "38.96100000",
	// }

	return ChipPayWithdrawCallbackReq{
		Sign:            "Pnm/E0715V3lpQzmiqV0wUTwTNhCRf8nXKXwBoCxD2hR6V66tadvqRvrQGvRfx9DniBrecZMYh2QQPhndfP6iZeOEZMHyzZeSvr5L3he63o8gp6Wd1JRQmRp+gAHmet9gfZbYIOxZP7++F2A+yX7Bcdeip6LMcIZDdCaBABq+io=",
		CompanyOrderNum: "202603260856400878",
		OtcOrderNum:     "4548687456245761_17745058799019",
		CancelReason:    "用户收款帐户屬於高風險银行 User’s receiving account belongs to a high-risk bank.",
		CoinAmount:      "35.84410000",
		CoinSign:        "usdt",
		OrderType:       "2",
		TradeStatus:     "0",
		TradeOrderTime:  "2026-03-26 14:18:00",
		UnitPrice:       "7.70000000",
		Total:           "276.00000000",
		SuccessAmount:   "0",
	}
}
