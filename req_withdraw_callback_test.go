package go_chipay

import (
	"encoding/json"
	"testing"
)

func TestWithdrawCallback(t *testing.T) {
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

	req := GenWdCallbackRequestDemo()
	body, err := json.Marshal(req)
	if err != nil {
		cli.logger.Errorf("Error:%s", err.Error())
		t.Fatal(err)
	}

	err = cli.WithdrawCallback(body, func(ChipPayWithdrawCallbackReq) error { return nil })
	if err != nil {
		cli.logger.Errorf("Error:%s", err.Error())
		t.Fatal(err)
	}
	cli.logger.Infof("resp:%+v\n", req)
}

func GenWdCallbackRequestDemo() ChipPayWithdrawCallbackReq {
	return ChipPayWithdrawCallbackReq{
		Sign:            "TYzuqA6gkS9wAzLA+MTx5/6TF89cH3JsTkC67WDi8u8NlGp5DuCLIltRzO/c8aG9h54dLYyxtFrrTsB9qgWmpXNMXqPX1PNeIoLS4D/l/jh6yIohTwpKvp5Giq7UffnMgyt6ha2OGc8kam6rilX+ZLi7CP6mcx/DYkO46c8b8q0=",
		CompanyOrderNum: "hafagafasfadfwerwer32",
		OtcOrderNum:     "12511234561_1592731510161",
		CoinAmount:      "100.00",
		CoinSign:        "usdt",
		OrderType:       "2",
		TradeStatus:     "1",
		TradeOrderTime:  "2020-07-15 18:46:04",
		UnitPrice:       "7.01",
		Total:           "701",
		SuccessAmount:   "100",
	}
}
