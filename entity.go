package go_chipay

type ChipPayInitParams struct {
	MerchantInfo `yaml:",inline" mapstructure:",squash"`

	DepositURL       string `json:"depositURL" mapstructure:"depositURL" config:"depositURL" yaml:"depositURL"`                         // 入金下单接口地址
	WithdrawURL      string `json:"withdrawURL" mapstructure:"withdrawURL" config:"withdrawURL" yaml:"withdrawURL"`                     // 出金下单接口地址
	QueryDepositURL  string `json:"queryDepositURL" mapstructure:"queryDepositURL" config:"queryDepositURL" yaml:"queryDepositURL"`     // 入金订单查询接口地址
	QueryWithdrawURL string `json:"queryWithdrawURL" mapstructure:"queryWithdrawURL" config:"queryWithdrawURL" yaml:"queryWithdrawURL"` // 出金订单查询接口地址
}

type MerchantInfo struct {
	CompanyID         int64  `json:"companyId" mapstructure:"companyId" config:"companyId" yaml:"companyId"`                                 // 商户ID
	PrivateKey        string `json:"privateKey" mapstructure:"privateKey" config:"privateKey" yaml:"privateKey"`                             // 商户私钥（Base64 DER PKCS#8 或 PEM）
	CallbackPublicKey string `json:"callbackPublicKey" mapstructure:"callbackPublicKey" config:"callbackPublicKey" yaml:"callbackPublicKey"` // ChipPay 回调公钥
	AsyncURL          string `json:"asyncUrl" mapstructure:"asyncUrl" config:"asyncUrl" yaml:"asyncUrl"`                                     // 异步通知地址 (商户接收回调通知的地址)
	SyncURL           string `json:"syncUrl" mapstructure:"syncUrl" config:"syncUrl" yaml:"syncUrl"`                                         // 同步返回地址 (用户完成或取消交易后返回至商户平台的地址)
}

type ChipPayDepositReq struct {
	CompanyID            int64   `json:"companyId"`                      // 商户ID
	CompanyOrderNum      string  `json:"companyOrderNum"`                // 商户订单号（唯一）
	AreaCode             string  `json:"areaCode,omitempty"`             // 区号，默认 86
	Phone                string  `json:"phone,omitempty"`                // 用户手机号
	AsyncURL             string  `json:"asyncUrl"`                       // 异步通知地址
	SyncURL              string  `json:"syncUrl"`                        // 同步跳转地址
	TotalAmount          *int64  `json:"totalAmount,omitempty"`          // 法币总金额（整数，和 CoinQuantity 二选一，优先 totalAmount）
	CoinQuantity         *string `json:"coinQuantity,omitempty"`         // USDT 数量（最多4位小数）
	Name                 string  `json:"name,omitempty"`                 // 真实姓名
	IdentityType         *int    `json:"identityType,omitempty"`         // 证件类型：1身份证，2护照
	Area                 string  `json:"area,omitempty"`                 // 国家地区（ISO3166-1 三位字母）
	Number               string  `json:"number,omitempty"`               // 证件号
	IdentityPictureFront string  `json:"identityPictureFront,omitempty"` // 证件人像面URL
	IdentityPictureBack  string  `json:"identityPictureBack,omitempty"`  // 证件无人像面URL
	AdditionalPicture    string  `json:"additionalPicture,omitempty"`    // 附加证件URL
	PaymentMethod        string  `json:"paymentMethod,omitempty"`        // 支付方式：1银行卡 2支付宝 3微信
	Sign                 string  `json:"sign,omitempty"`                 // 请求签名（SDK自动计算）
}

type ChipPayDepositRsp struct {
	Code    int                `json:"code"`    // 状态码
	Msg     string             `json:"msg"`     // 提示信息
	Data    ChipPayDepositData `json:"data"`    // 业务数据
	Success bool               `json:"success"` // 是否成功
}

type ChipPayDepositData struct {
	Link          string `json:"link"`          // 收银台链接
	IntentOrderNo string `json:"intentOrderNo"` // ChipPay 意向订单号
}

type ChipPayDepositQueryReq struct {
	CompanyID       int64  `json:"companyId"`       // 商户ID
	CompanyOrderNum string `json:"companyOrderNum"` // 商户订单号
	Sign            string `json:"sign,omitempty"`  // 请求签名（SDK自动计算）
}

type ChipPayDepositQueryRsp struct {
	Code    int                     `json:"code"`    // 状态码
	Msg     string                  `json:"msg"`     // 提示信息
	Data    ChipPayDepositQueryData `json:"data"`    // 业务数据
	Success bool                    `json:"success"` // 是否成功
}

type ChipPayDepositQueryData struct {
	IntentOrderNo   int64  `json:"intentOrderNo"`    // ChipPay 意向订单号
	CompanyOrderNum string `json:"companyOrderNum"`  // 商户订单号
	TradeOrderTime  string `json:"tradeOrderTime"`   // 订单创建时间
	SuccessAmount   string `json:"successAmount"`    // 实际到账币数量
	UnitPrice       string `json:"unitPrice"`        // 币种单价
	Total           string `json:"total"`            // 法币总额
	CoinAmount      string `json:"coinAmount"`       // 初始币数量
	CoinSign        string `json:"coinSign"`         // 币种标识（示例5=USDT）
	TradeStatus     string `json:"tradeStatus"`      // 交易状态：0失败 1成功 2处理中
	ErrorOrderTotal string `json:"errorOrderTotal"`  // 异常处理后的法币金额
	ErrorOrderAmt   string `json:"errorOrderAmount"` // 异常处理后的币到账数量
	Sign            string `json:"sign"`             // 响应签名
}

type ChipPayDepositCallbackReq struct {
	CompanyOrderNum string `json:"companyOrderNum"` // 商户订单号
	IntentOrderNo   string `json:"intentOrderNo"`   // ChipPay 意向订单号
	CoinAmount      string `json:"coinAmount"`      // 初始币数量
	CoinSign        string `json:"coinSign"`        // 币种标识
	TradeStatus     string `json:"tradeStatus"`     // 交易状态：0失败 1成功
	CancelReason    string `json:"cancelReason"`    // 取消原因（失败时返回）
	TradeOrderTime  string `json:"tradeOrderTime"`  // 订单时间
	UnitPrice       string `json:"unitPrice"`       // 单价
	Total           string `json:"total"`           // 法币实际到账金额
	SuccessAmount   string `json:"successAmount"`   // 币种到账数量
	Sign            string `json:"sign"`            // 回调签名
}

type ChipPayWithdrawReq struct {
	CompanyID        int64                  `json:"companyId"`                  // 商户ID
	Kyc              string                 `json:"kyc"`                        // 用户验证级别（默认2）
	Username         string                 `json:"username"`                   // 真实姓名
	AreaCode         string                 `json:"areaCode,omitempty"`         // 国际区号（默认86）
	Phone            string                 `json:"phone,omitempty"`            // 手机号（phone/email 二选一）
	Email            string                 `json:"email,omitempty"`            // 邮箱（仅vnd场景可用，phone/email 二选一）
	OrderType        int                    `json:"orderType"`                  // 订单类型（1买单，2卖单；出金固定2）
	IDCardType       *int                   `json:"idCardType,omitempty"`       // 证件类型（1身份证 2护照 3其他）
	IDCardNum        string                 `json:"idCardNum,omitempty"`        // 证件号码
	PayCardNo        string                 `json:"payCardNo,omitempty"`        // 收款银行卡号（卖单必填，仅数字）
	PayCardBank      string                 `json:"payCardBank,omitempty"`      // 开户银行（卖单必填）
	PayCardBranch    string                 `json:"payCardBranch,omitempty"`    // 开户支行
	CompanyOrderNum  string                 `json:"companyOrderNum"`            // 商户订单号
	CoinSign         string                 `json:"coinSign"`                   // 数字货币标识（USDT）
	PayCoinSign      string                 `json:"payCoinSign"`                // 法币币别（cny/vnd）
	CoinAmount       string                 `json:"coinAmount,omitempty"`       // 下单USDT数量（和total二选一，total优先）
	Total            string                 `json:"total,omitempty"`            // 法币金额（卖单可小数）
	OrderPayChannel  *int                   `json:"orderPayChannel,omitempty"`  // 支付渠道（cny卖单固定3银行卡）
	DisplayUnitPrice string                 `json:"displayUnitPrice,omitempty"` // 自定义展示单价（最多4位小数）
	OrderTime        string                 `json:"orderTime"`                  // 订单时间戳（毫秒，5分钟内有效）
	SyncURL          string                 `json:"syncUrl,omitempty"`          // 同步返回地址
	AsyncURL         string                 `json:"asyncUrl,omitempty"`         // 异步通知地址
	Extra            map[string]interface{} `json:"extra,omitempty"`            // 扩展字段（用于未来新增参数）
	Sign             string                 `json:"sign,omitempty"`             // 请求签名（SDK自动计算）
}

type ChipPayWithdrawRsp struct {
	Code    int                 `json:"code"`    // 状态码
	Msg     string              `json:"msg"`     // 提示信息
	Data    ChipPayWithdrawData `json:"data"`    // 业务数据
	Success bool                `json:"success"` // 是否成功
}

type ChipPayWithdrawCallbackReq struct {
	CompanyOrderNum string `json:"companyOrderNum"` // 商户订单号
	OtcOrderNum     string `json:"otcOrderNum"`     // ChipPay 平台订单号
	CoinAmount      string `json:"coinAmount"`      // 订单初始币数量
	CoinSign        string `json:"coinSign"`        // 币种标识（usdt）
	OrderType       string `json:"orderType"`       // 订单类型（1买单，2卖单）
	TradeStatus     string `json:"tradeStatus"`     // 交易状态（0失败 1成功 2批量卖单生成失败 4商户手动取消）
	CancelReason    string `json:"cancelReason"`    // 取消原因（失败时返回）
	TradeOrderTime  string `json:"tradeOrderTime"`  // 订单交易时间（北京时间）
	UnitPrice       string `json:"unitPrice"`       // 币种单价
	Total           string `json:"total"`           // 法币实际到账金额
	SuccessAmount   string `json:"successAmount"`   // 实际到账币数量
	Sign            string `json:"sign"`            // 回调签名
}

type ChipPayWithdrawData struct {
	Link    string `json:"link"`    // 收银台链接
	OrderNo string `json:"orderNo"` // ChipPay 平台订单号
}

type ChipPayCallbackRsp struct {
	Code int    `json:"code"` // 回调响应码（200表示成功）
	Msg  string `json:"msg"`  // 回调响应信息
}

// 入金回调返回
type DepositCallbackResponse struct {
	Code    int      `json:"code"`
	Meg     string   `json:"msg"`
	Data    BackBody `json:"data"`
	Success bool     `json:"success"`
}

// Id            string `json:"id"`
// OrderNo       string `json:"orderNo"`
// NotifyUrl     string `json:"notifyUrl"`
// NotifyReq     string `json:"notifyReq"`
// NotifyCount   string `json:"notifyCount"`
// NotifyStatus  string `json:"notifyStatus"`
// ErrorResponse string `json:"errorResponse"`
// UpdateTime    string `json:"updateTime"`
// SignType      string `json:"signType"`
type BackBody struct {
	IntentOrderNo   string `json:"intentOrderNo"`   // ChipPay 意向订单号
	CompanyOrderNum string `json:"companyOrderNum"` // 商户订单号
}
