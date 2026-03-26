package go_chipay

import (
	"testing"
)

func TestWithdraw(t *testing.T) {
	vLog := VLog{}
	cli := NewClient(vLog, &ChipPayInitParams{
		MerchantInfo: MerchantInfo{
			CompanyID:         COMPANY_ID,
			PrivateKey:        PRIVATE_KEY,
			CallbackPublicKey: CALLBACK_PUBLIC_KEY,
			SyncURL:           SYNC_URL,
			WithdrawAsyncUrl:  WITHDRAW_ASYNC_URL,
		},
		DepositURL:       DEPOSIT_URL,
		WithdrawURL:      WITHDRAW_URL,
		QueryDepositURL:  QUERY_DEPOSIT_URL,
		QueryWithdrawURL: QUERY_WITHDRAW_URL,
	})

	resp, err := cli.WithdrawReq(GenWithdrawRequestDemo())
	if err != nil {
		cli.logger.Errorf("err:%s\n", err.Error())
		return
	}
	cli.logger.Infof("resp:%+v\n", resp)
}

func GenWithdrawRequestDemo() ChipPayWithdrawReq {
	return ChipPayWithdrawReq{
		CompanyOrderNum: "2029884682646034",
		Username:        "赫敏·珍珍·格兰杰",
		Phone:           "5300231651",
		PayCardNo:       "622848202009358999",
		PayCardBank:     "中国银行",
		PayCardBranch:   "上海",
		CoinSign:        "USDT",
		PayCoinSign:     "cny", // 法币币别，须传小写英文(cny，vnd)
		Total:           "300", // 用户付款的法币总金额(快捷买单只能传整数，快捷卖单不限)
	}
}
