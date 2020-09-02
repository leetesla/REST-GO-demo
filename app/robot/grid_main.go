package robot

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"REST-GO-demo/app/models"
)

//GridBuy 网格策略,只卖币,用户一段时间后手动买币
type GridBuy struct {
	Id   int
	User string
	//火币子账户id
	AccountId string
	//模式 0,正常模式 1,自动买手动卖 2,自动卖手动买
	Model int64
	//例如 btc
	BaseCurrency string
	//计量币 例如 usdt
	QuoteCurrency string

	//价格精度 小数点后几位
	PricePrecision uint64
	//数量精度
	AmountPrecision uint64
	//货币交易所api
	HuoBi *models.HuoBiEx
	//执行上下文
	Datas []*Ctxt
	//停止运行标记
	Istop bool
	//0 stop 1 running
	Status uint64
}

//ctxt 策略数据
type Ctxt struct {
	BuyMount float64 //买的量 usdt 计量
	BuyPrice float64 //买的价格
	BuyOrder string  //买的订单

	//SellMount float64 //卖的量 btc 计量
	SellPrice float64 //卖的价格
	SellOrder string  //卖的订单

	TotalBaseCurrency float64 // 到此需要买的币(btc) 用于初始化
	//TotalQuoteCurrency float64 //
}

//
func NewGridBuy(ID int, Model int64, AccountId, BaseCurrency, QuoteCurrency, Datas, user string, AccessKay, SecretKet string) *GridBuy {
	m := new(GridBuy)
	m.Id = ID
	m.Model = Model
	err := json.Unmarshal([]byte(Datas), &m.Datas)
	if err != nil {
		panic(err.Error())
	}
	m.User = user
	m.AccountId = AccountId
	m.BaseCurrency = BaseCurrency
	m.QuoteCurrency = QuoteCurrency
	m.HuoBi = models.NewHuoBiEx(AccessKay, SecretKet)

	return m
}

//Register 注册服务方便外界调用
func (m *GridBuy) Register() {
	if err := Register(m.Id, m); err != nil {
		log.Println(err.Error())
		Stop(m.Id)
		Register(m.Id, m)
	}
}

//Start 启动服务
func (m *GridBuy) Start() {
	err := m.Init()
	if err != nil {
		log.Println(err.Error())
		return
	}
	m.Register()
	// 标记为running
	m.Status = 1

	if m.Model == 1 {
		m.work_classical()
	} else if m.Model == 2 {
		m.work_buy()
	} else if m.Model == 3 {
		m.work_sell()
	}
	m.Status = 0
}

//Init 获取交易对精度
func (m *GridBuy) Init() error {
	sbls := m.HuoBi.GetSymbols()
	//log.Println(sbls)
	if sbls.Status != "ok" {
		return errors.New(sbls.ErrMsg)
	}
	for i := 0; i < len(sbls.Data); i++ {
		if sbls.Data[i].BaseCurrency == m.BaseCurrency &&
			sbls.Data[i].QuoteCurrency == m.QuoteCurrency {
			m.PricePrecision = uint64(sbls.Data[i].PricePrecision)
			m.AmountPrecision = uint64(sbls.Data[i].AmountPrecision)
			return nil
		}
	}
	return nil
}

//sleep 睡眠 单位s
func (m *GridBuy) sleepOrStop(ts uint64) error {
	for i := 0; uint64(i) < ts; i++ {
		time.Sleep(time.Second)
		if m.Istop == true {
			return errors.New("stop sg")
		}
	}
	return nil
}

// Stop 停止服务,外界调用
func (m *GridBuy) Stop() {
	m.Istop = true
	for {
		if m.Status == 0 {
			return
		}
		time.Sleep(time.Second)
	}
}

//Action 外界调用
//只买模式触发卖
//只卖模式触发买
func (m *GridBuy) Action() error {
	if m.Model == 1 {
		return m.buy_action()
	} else if m.Model == 2 {
		return m.sell_action()
	}
	return nil
}
