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
			AsyncURL:          ASYNC_URL,
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
		CompanyOrderNum: "37028219790919996X",
		Username:        "赫敏·珍珍·格兰杰",
		Phone:           "5300231651",
		PayCardNo:       "622848202009358999",
		PayCardBank:     "中国银行",
		PayCardBranch:   "上海",
		CoinSign:        "USDT",
		PayCoinSign:     "cny",
		Total:           "300",
	}
}
