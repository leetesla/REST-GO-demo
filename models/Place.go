package models

type PlaceRequestParams struct {
	AccountID string `json:"account-id"` // 账户ID
	Amount    string `json:"amount"`     // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
	Price     string `json:"price"`      // 下单价格, 市价单不传该参数
	Source    string `json:"source"`     // 订单来源, api: API调用, margin-api: 借贷资产交易
	Symbol    string `json:"symbol"`     // 交易对, btcusdt, bccbtc......
	Type      string `json:"type"`       // 订单类型, buy-market: 市价买, sell-market: 市价卖, buy-limit: 限价买, sell-limit: 限价卖
}

type PlaceReturn struct {
	Status  string `json:"status"`
	Data    string `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type OrderReturn struct {
	Status string `json:"status"`
	Data   Order  `json:"data"`
}

type Order struct {
	Id              uint64  `json:"id"`
	Symbol          string  `json:"symbol"`
	Accountid       uint64  `json:"account-id"`
	Amount_f        float64 `json:"-"`
	Amount          string  `json:"amount"` //订单数量
	Price_f         float64 `json:"-"`
	Price           string  `json:"price"`
	Createdat       uint64  `json:"created-at"`
	Type            string  `json:"type"`
	Fieldamount     string  `json:"field-amount"`      //已成交数量
	Fieldcashamount string  `json:"field-cash-amount"` //已成交总金额
	Fieldfees       string  `json:"field-fees"`        //已成交手续费（买入为币，卖出为钱）
	Finishedat      uint64  `json:"finished-at"`       //订单变为终结态的时间，不是成交时间，包含“已撤单”状态
	Userid          uint64  `json:"user-id"`
	Source          string  `json:"source"`
	//submitted 已提交, partial-filled 部分成交, partial-canceled 部分成交撤销, filled 完全成交, canceled 已撤销， created
	State      string `json:"state"`
	Canceledat uint64 `json:"canceled-at"`
	Exchange   string `json:"exchange"`
	Batch      string `json:"batch"`
}
