package go_chipay

import (
	"testing"
)

func TestDeposit(t *testing.T) {
	vLog := VLog{}
	cli := NewClient(vLog, &ChipPayInitParams{
		MerchantInfo: MerchantInfo{
			CompanyID:         COMPANY_ID,
			PrivateKey:        PRIVATE_KEY,
			CallbackPublicKey: CALLBACK_PUBLIC_KEY,
			SyncURL:           SYNC_URL,
			DepositAsyncUrl:   DEPOSIT_ASYNC_URL,
		},
		DepositURL:       DEPOSIT_URL,
		WithdrawURL:      WITHDRAW_URL,
		QueryDepositURL:  QUERY_DEPOSIT_URL,
		QueryWithdrawURL: QUERY_WITHDRAW_URL,
	})

	resp, err := cli.Deposit(GenDepositRequestDemo())
	if err != nil {
		cli.logger.Errorf("err:%s\n", err.Error())
		return
	}
	cli.logger.Infof("resp:%+v\n", resp)
}

// areaCode: 90
// phone: 5300231651
// name: 赫敏·珍珍·格兰杰
// number: 37028219790919996X
func GenDepositRequestDemo() ChipPayDepositReq {
	amount := int64(300)
	return ChipPayDepositReq{
		CompanyOrderNum: "202666335563990",
		Phone:           "5300231651",
		TotalAmount:     &amount,
		Name:            "赫敏·珍珍·格兰杰",
		Number:          "37028219790919996X",
	}
}
