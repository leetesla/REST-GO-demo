package robot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	appmodels "REST-GO-demo/app/models"
	"REST-GO-demo/models"
)

//Now 获取服务当前时间戳
func (m *GridBuy) Now() (uint64, error) {
	ts := m.HuoBi.GetTimestamp()
	if ts.Status != "ok" || uint64(ts.Data/1000) == 0 {
		return 0, errors.New(ts.ErrMsg)
	}
	return uint64(ts.Data / 1000), nil
}

//BalanceOf 查询账户币额度
func (m *GridBuy) BalanceOf(coin string) (ret float64, err error) {
	//Type := "trade"
	balance := m.HuoBi.GetAccountBalance(m.AccountId)
	if balance.Status != "ok" {
		return 0, errors.New(balance.ErrMsg)
	}

	for i := 0; i < len(balance.Data.List); i++ {
		if balance.Data.List[i].Currency == coin {
			balance, _ := strconv.ParseFloat(balance.Data.List[i].Balance, 64)
			ret += balance
		}
	}
	return ret, nil
}

//BuyPrice 当前买价
func (m *GridBuy) BuyPrice() (float64, error) {
	ticker := m.HuoBi.GetTicker(m.BaseCurrency + m.QuoteCurrency)
	if ticker.Status != "ok" {

		return 0, errors.New(ticker.ErrMsg)
	}
	return ticker.Tick.Ask[0], nil
}

//SellPrice 当前卖价
func (m *GridBuy) SellPrice() (float64, error) {
	ticker := m.HuoBi.GetTicker(m.BaseCurrency + m.QuoteCurrency)
	if ticker.Status != "ok" {
		return 0, errors.New(ticker.ErrMsg)
	}
	return ticker.Tick.Bid[0], nil
}

//BuyMarket 市价买币
func (m *GridBuy) BuyMarket(_usdt float64) error {
	resq := models.PlaceRequestParams{
		AccountID: m.AccountId,
		Amount:    m.String(_usdt, m.AmountPrecision),
		Price:     "",
		Source:    "api",
		Symbol:    m.BaseCurrency + m.QuoteCurrency,
		Type:      "buy-market"}
	log.Println(m.AmountPrecision, resq)
	ret := m.HuoBi.Place(resq)
	if ret.Status != "ok" {
		return errors.New(ret.ErrMsg)
	}
	return new(appmodels.Orders).Add(ret.Data, m.Id)
}

//BuyOrder 限价买币订单
func (m *GridBuy) BuyOrder(_btc float64, _price float64) (string, error) {
	resq := models.PlaceRequestParams{
		AccountID: m.AccountId,
		Amount:    m.String(_btc, m.AmountPrecision),
		Price:     m.String(_price, m.PricePrecision), //m.String(_price, m.PricePrecision),
		Source:    "api",
		Symbol:    m.BaseCurrency + m.QuoteCurrency,
		Type:      "buy-limit"}
	//log.Println(resq)
	ret := m.HuoBi.Place(resq)
	if ret.Status != "ok" {
		return "", errors.New(ret.ErrMsg)
	}
	new(appmodels.Orders).Add(ret.Data, m.Id)
	return ret.Data, nil
}

//SellMarket 市价卖币
func (m *GridBuy) SellMarket(_btc float64) error {
	//log.Print("SellMarket", _btc, ",", m.AmountPrecision, m.String(_btc, m.AmountPrecision))
	resq := models.PlaceRequestParams{
		AccountID: m.AccountId,
		Amount:    m.String(_btc, m.AmountPrecision),
		Price:     "",
		Source:    "api",
		Symbol:    m.BaseCurrency + m.QuoteCurrency,
		Type:      "sell-market"}
	ret := m.HuoBi.Place(resq)
	if ret.Status != "ok" {
		return errors.New(ret.ErrMsg)
	}
	return new(appmodels.Orders).Add(ret.Data, m.Id)
}

//SellOrder 限价卖币订单
func (m *GridBuy) SellOrder(_btc float64, _price float64) (string, error) {
	resq := models.PlaceRequestParams{
		AccountID: m.AccountId,
		Amount:    m.String(_btc, m.AmountPrecision),
		Price:     m.String(_price, m.PricePrecision),
		Source:    "api",
		Symbol:    m.BaseCurrency + m.QuoteCurrency,
		Type:      "sell-limit"}
	ret := m.HuoBi.Place(resq)
	if ret.Status != "ok" {
		return "", errors.New(ret.ErrMsg)
	}
	new(appmodels.Orders).Add(ret.Data, m.Id)
	return ret.Data, nil
}

//CancelOrder 取消订单
func (m *GridBuy) CancelOrder(orderid string) error {
	ret := m.HuoBi.SubmitCancel(orderid)
	if ret.Status != "ok" {
		return errors.New(ret.ErrMsg)
	}
	return nil
}

//GetOrder 获取订单详情
func (m *GridBuy) GetOrder(orderid string) (models.OrderReturn, error) {
	order := m.HuoBi.GetOrder(orderid)
	if order.Status != "ok" {
		return order, errors.New(fmt.Sprint("cont get order:", orderid))
	}
	order.Data.Amount_f, _ = strconv.ParseFloat(order.Data.Amount, 64)
	order.Data.Price_f, _ = strconv.ParseFloat(order.Data.Price, 64)
	return order, nil
}

//String float64转string
//params.decimal 保留小数位数
func (m *GridBuy) String(f float64, decimal uint64) string {
	if decimal == 0 {
		return fmt.Sprintf("%v", uint64(f))
	} else if decimal <= 10 {
		return fmt.Sprintf(fmt.Sprintf("%%.%vf", decimal), f)
	} else {
		return fmt.Sprintf("%.10f", f)
	}
}

func (m *GridBuy) Save() error {
	db := new(appmodels.Stracy)
	data, _ := json.Marshal(m.Datas)
	return db.Save(m.Id, string(data))
}
