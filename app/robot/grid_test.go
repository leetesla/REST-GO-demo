package robot

import "testing"

func TestGridInit(t *testing.T) {
	grid := NewGridBuy(1, 1, "", "btc", "usdt", "[]","","","")
	err := grid.Init()
	if err != nil {
		t.Error("获取精度失败:", grid.PricePrecision, grid.AmountPrecision)
	}
	t.Log("获取精度:", grid.PricePrecision, grid.AmountPrecision)
	return
}
func TestGridNow(t *testing.T) {
	grid := NewGridBuy(1, 1, "", "btc", "usdt", "[]","","","")
	now, err := grid.Now()
	if err != nil {
		t.Error("获取服务器时间戳失败:", err.Error())
	}
	t.Log("获取服务器时间:", now)
}
func TestGridGetAccount(t *testing.T) {
	grid := NewGridBuy(1, 1, "8101421", "btc", "usdt", "[]","","","")
	accounts := grid.HuoBi.GetAccounts()
	if accounts.Status != "ok" {
		t.Error("获取账户失败:", accounts.ErrMsg)
	}
	t.Log("accounts", accounts)
}

func TestGridBalance(t *testing.T) {
	grid := NewGridBuy(1, 1, "8101421", "btc", "usdt", "[]","","","")
	balance, err := grid.BalanceOf("usdt")
	if err != nil {
		t.Error("获取余额失败:", balance)
	}
	t.Log("余额", balance)
}
func TestGridPrice(t *testing.T) {
	grid := NewGridBuy(1, 1, "8101421", "btc", "usdt", "[]","","","")
	BuyPrice, err := grid.BuyPrice()
	if err != nil {
		t.Error("获取价格失败:", BuyPrice)
	}
	SellPrice, err := grid.SellPrice()
	if err != nil {
		t.Error("获取价格失败:", SellPrice)
	}
	t.Log("买价", BuyPrice, "卖价", SellPrice)
}

func TestGridString(t *testing.T) {
	grid := NewGridBuy(1, 1, "8101421", "btc", "usdt", "[]","","","")

	t.Log(grid.String(1.1234567891, 0),
		grid.String(1.1234567891, 1),
		grid.String(1.1234567891, 2),
		grid.String(1.1234567891, 3),
		grid.String(1.1234567891, 4),
		grid.String(1.1234567891, 5),
		grid.String(1.1234567891, 6),
		grid.String(1.1234567891, 8),
		grid.String(1.1234567891, 9),
		grid.String(1.1234567891, 10),
		grid.String(1.1234567891, 11),
		grid.String(1.1234567891, 12))
}

func TestGridMarket(t *testing.T) {
	return
	grid := NewGridBuy(1, 1, "8101421", "btc", "usdt", "[]","","","")
	err := grid.Init()
	if err != nil {
		t.Error("获取精度失败:", grid.PricePrecision, grid.AmountPrecision)
	}
	err = grid.BuyMarket(10)
	if err != nil {
		t.Error("市价买失败:", err.Error())
		return
	}
	btc, err := grid.BalanceOf("btc")
	if err != nil {
		t.Error("获取余额失败:", err.Error())
	}
	t.Log("balance of btc:", btc)
	err = grid.SellMarket(btc)
	if err != nil {
		t.Error("市价卖失败:", err.Error())
	}
}
func TestGridBuyOrder(t *testing.T) {
	grid := NewGridBuy(1, 1, "8101421", "btc", "usdt", "[]","","","")
	err := grid.Init()
	if err != nil {
		t.Error("获取精度失败:", grid.PricePrecision, grid.AmountPrecision)
	}
	price, err := grid.BuyPrice()
	if err != nil {
		t.Error(err.Error())
	}
	orderid, err := grid.BuyOrder(0.001, price*0.99)
	if err != nil {
		t.Error(err.Error())
		return
	}
	//t.Log(order)
	orderReturn, err := grid.GetOrder(orderid)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(orderReturn)
	err = grid.CancelOrder(orderid)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGridSellOrder(t *testing.T) {
	grid := NewGridBuy(1, 1, "8101421", "btc", "usdt", "[]","","","")
	err := grid.Init()
	if err != nil {
		t.Error(err.Error())
	}
	price, err := grid.SellPrice()
	if err != nil {
		t.Error(err.Error())
	}
	orderid, err := grid.SellOrder(0.0001, price*1.1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	//t.Log(order)
	err = grid.CancelOrder(orderid)
	if err != nil {
		t.Error(err.Error())
	}
}
