package go_chipay

import (
	"testing"
)

func newUSDTDepositClient() *Client {
	vLog := VLog{}
	return NewClient(vLog, &ChipPayInitParams{
		MerchantInfo: MerchantInfo{
			CompanyID:         COMPANY_ID,
			PrivateKey:        PRIVATE_KEY,
			CallbackPublicKey: CALLBACK_PUBLIC_KEY,
			DepositAsyncUrl:   DEPOSIT_ASYNC_URL,
			SyncURL:           SYNC_URL,
		},
		DepositURL:       DEPOSIT_URL,
		WithdrawURL:      WITHDRAW_URL,
		QueryDepositURL:  QUERY_DEPOSIT_URL,
		QueryWithdrawURL: QUERY_WITHDRAW_URL,
		USDTDepositURL:   USDT_DEPOSIT_URL,
	})
}

// TestUSDTDeposit 测试 USDT 轮换地址模式下单（1.1.3.1）。
// rechargeCoin=5 表示 USDT，rechargeAmount 不传表示不限定金额。
// 不传 CoinType 时收银台默认展示所有支持链的充值地址。
func TestUSDTDeposit(t *testing.T) {
	cli := newUSDTDepositClient()

	resp, err := cli.USDTDeposit(ChipPayUSDTDepositReq{
		CompanyOrderNo: "USDT20260008",  // 替换为真实唯一订单号
		RechargeCoin:   "5",             // USDT
		RechargeAmount: float64ptr(100), // 可选：指定存币数量（double）
		CoinType:       "trc20",         // 可选：指定收银台仅显示该链类型，如 "erc20"、"trc20"、"bep20"、"spl"
	})
	if err != nil {
		cli.logger.Errorf("USDTDeposit err: %s", err.Error())
		return
	}
	cli.logger.Infof("USDTDeposit resp: %+v", resp)
	if resp.Data.Link != "" {
		cli.logger.Infof("收银台链接: %s", resp.Data.Link)
	}
	for _, addr := range resp.Data.OrderAddressesDTO.Addresses {
		cli.logger.Infof("链=%s 协议=%s 地址=%s", addr.Mainnet, addr.Protocol, addr.Address)
	}
}

// TestUSDTDepositWithCoinType 测试 USDT 轮换地址模式下单并指定链类型（1.1.3.1）。
// 设置 CoinType 后，SDK 将在返回的收银台 Link 末尾自动追加 &cointype=<value>，
// 使收银台仅显示对应链类型的充值地址。
func TestUSDTDepositWithCoinType(t *testing.T) {
	cli := newUSDTDepositClient()

	resp, err := cli.USDTDeposit(ChipPayUSDTDepositReq{
		CompanyOrderNo: "USDT20260004", // 替换为真实唯一订单号
		RechargeCoin:   "5",            // USDT
		CoinType:       "trc20",        // 仅在收银台显示 TRC20 链的充值地址
	})
	if err != nil {
		cli.logger.Errorf("USDTDepositWithCoinType err: %s", err.Error())
		return
	}
	cli.logger.Infof("USDTDepositWithCoinType resp: %+v", resp)
	if resp.Data.Link != "" {
		// Link 末尾已自动追加 &cointype=trc20
		cli.logger.Infof("收银台链接（含链类型限定）: %s", resp.Data.Link)
	}
	for _, addr := range resp.Data.OrderAddressesDTO.Addresses {
		cli.logger.Infof("链=%s 协议=%s 地址=%s", addr.Mainnet, addr.Protocol, addr.Address)
	}
}

// TestUSDTDepositCallback 测试 USDT 轮换地址模式回调验签（1.1.3.2）。
// 将真实回调 JSON 参数填入 req 进行验签验证。
func TestUSDTDepositCallback(t *testing.T) {
	cli := newUSDTDepositClient()

	// 真实回调数据（来自 2026-07-24 15:56:35 的通道回调日志）
	req := ChipPayUSDTDepositCallbackReq{
		CompanyOrderNo:          "USDT20260004",
		AccumulativeTotalAmount: "0.00100000",
		RechargeActualCoin:      "5",
		RechargeActualAmount:    "0.00100000",
		RechargeAccountTime:     "2026-07-24 15:56:35",
		RechargeDetailOrderNo:   "RECHARGE_DETAIL_4548687456245761_1784879794721",
		RechargeAddress:         "TBQcqee9tKBXXvaXiq7ERhdrJc2DWf6PM7",
		TxHash:                  "-",
		FromAddress:             "86-16451231217",
		Sign:                    "U4tKzw8ZBDPIcVRoQgJmuCPEipz+mnmcbq4S8nheZTveRvC3Ci34q03oy0N/DwGlFMs/0i6yCywuYSs9HKpKCBQgrPKc3dTzq5fhrK3cut/vLJeFDkhM1FjUVbqEPvhDDVIDqJcI9m+H8FkqKR5pizZMjDrD2i+UoReqRqNXdCI=",
	}

	err := cli.USDTDepositCallback(req, func(r ChipPayUSDTDepositCallbackReq) error {
		cli.logger.Infof("USDTDepositCallback 验签通过, companyOrderNo=%s, rechargeActualAmount=%s, rechargeDetailOrderNo=%s, rechargeAddress=%s, txHash=%s",
			r.CompanyOrderNo, r.RechargeActualAmount, r.RechargeDetailOrderNo, r.RechargeAddress, r.TxHash)
		return nil
	})
	if err != nil {
		cli.logger.Errorf("USDTDepositCallback err: %s", err.Error())
		t.Fatal(err)
	}
}

// float64ptr 返回 float64 指针，用于填写可选的 RechargeAmount 字段（double 类型）
func float64ptr(f float64) *float64 {
	return &f
}
