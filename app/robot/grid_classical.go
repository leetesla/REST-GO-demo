package robot

import (
	"log"

	"github.com/onethefour/REST-GO-demo/app/utils"
)

// work_classical 经典网格策略
func (m *GridBuy) work_classical() {
	var affacted bool = true
	var err error
	if err = m.classica_init(); err != nil {
		log.Println(err.Error())
		return
	}

	for {
		if affacted {
			if err = m.Save(); err != nil {
				log.Println(err.Error())
			}
		}

		if err = m.sleepOrStop(10); err != nil {
			m.Status = 0
			m.Istop = false
			log.Println("serve stop")
			return
		}

		if affacted, err = m.classical_clear_order(); affacted || err != nil {
			if err != nil {
				log.Println(err.Error())
			}
			continue
		}
		if affacted, err = m.classica_new_order(); affacted || err != nil {
			if err != nil {
				log.Println(err.Error())
			}
			continue
		}
	}
}

//classica_init 初始检查
func (m *GridBuy) classica_init() (err error) {
	//计算各个高度需要买的btc
	for i := 0; i < len(m.Datas); i++ {
		if i == 0 {
			m.Datas[i].TotalBaseCurrency = m.Datas[i].BuyMount / m.Datas[i].BuyPrice
		} else {
			m.Datas[i].TotalBaseCurrency = m.Datas[i-1].TotalBaseCurrency + m.Datas[i].BuyMount/m.Datas[i].BuyPrice
		}
	}
	var buyPrice, sellPrice, buyNeedBtc, sellNeedBtc, nowBtc float64
	if buyPrice, err = m.BuyPrice(); err != nil {
		log.Println(err.Error())
		return err
	}
	if sellPrice, err = m.SellPrice(); err != nil {
		log.Println(err.Error())
		return err
	}
	//当前价格所在的高度至少需要买多少btc
	for i := 0; i < len(m.Datas); i++ {
		if m.Datas[i].BuyPrice > buyPrice {
			buyNeedBtc = m.Datas[i].TotalBaseCurrency
		} else {
			break
		}
	}
	//当前价格所在的高度至多需要买多少btc
	for i := 0; i < len(m.Datas); i++ {
		if m.Datas[i].SellPrice > sellPrice {
			sellNeedBtc = m.Datas[i].TotalBaseCurrency
		} else {
			break
		}
	}
	//当前账户上btc
	if nowBtc, err = m.BalanceOf(m.BaseCurrency); err != nil {
		log.Println(err.Error())
		return err
	}
	//log.Panicln("needBtc", needBtc, "nowBtc", nowBtc, m.Datas[0].TotalBaseCurrency)
	//补些btc,一般是首次启动初始化
	if nowBtc < buyNeedBtc-(0.5*m.Datas[0].TotalBaseCurrency) {
		if err = m.BuyMarket((buyNeedBtc - nowBtc) * buyPrice); err != nil {
			return err
		}

		if nowBtc, err = m.BalanceOf(m.BaseCurrency); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	//卖掉多出的btc
	if nowBtc > sellNeedBtc+(0.5*m.Datas[0].TotalBaseCurrency) {
		if err = m.SellMarket(nowBtc - sellNeedBtc); err != nil {
			log.Println(err.Error())
			return err
		}
		if nowBtc, err = m.BalanceOf(m.BaseCurrency); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	//策略改动后,之前的订单需要重新下单
	for i := 0; i < len(m.Datas); i++ {
		orderId := m.Datas[i].BuyOrder
		orderPrice := m.Datas[i].BuyPrice
		orderAmount := m.Datas[i].BuyMount / m.Datas[i].BuyPrice
		//m.Datas[i].BuyOrder = ""
		if orderId == "" {
			orderId = m.Datas[i].SellOrder
			orderPrice = m.Datas[i].SellPrice
			//m.Datas[i].SellOrder = ""
		}
		if orderId != "" {
			orderReturn, err := m.GetOrder(orderId)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			if orderReturn.Data.Price_f > orderPrice*1.005 || orderReturn.Data.Price_f < orderPrice*0.996 {
				m.CancelOrder(orderId)
				continue
			}
			if orderReturn.Data.Amount_f > orderAmount*1.005 || orderReturn.Data.Amount_f < orderAmount*0.996 {
				m.CancelOrder(orderId)
				continue
			}
		}
	}
	return nil
}

//classical_clear_order 清除已经终结的订单
func (m *GridBuy) classical_clear_order() (affected bool, err error) {
	for i := 0; i < len(m.Datas); i++ {
		orderId := m.Datas[i].BuyOrder
		if orderId != "" {
			orderReturn, err := m.GetOrder(orderId)
			if err != nil {
				return affected, err
			}
			if orderReturn.Data.Createdat > 0 {
				m.Datas[i].BuyOrder = ""
				affected = true
			}
		}
		orderId = m.Datas[i].SellOrder
		if orderId != "" {
			orderReturn, err := m.GetOrder(orderId)
			if err != nil {
				return affected, err
			}
			if orderReturn.Data.Finishedat > 0 {
				m.Datas[i].SellOrder = ""
				affected = true
			}
		}
	}
	return
}

//下买单和卖单
func (m *GridBuy) classica_new_order() (affected bool, err error) {
	var buyPrice, sellPrice, buyNeedBtc, sellNeedBtc, nowBtc float64
	//需要下卖单的的最低高度
	sellIndex := -1
	//需要下买单的最高高度
	buyIndex := -1
	//当前价格
	if buyPrice, err = m.BuyPrice(); err != nil {
		log.Println(err.Error())
		return affected, err
	}
	//当前价格
	if sellPrice, err = m.SellPrice(); err != nil {
		log.Println(err.Error())
		return affected, err
	}
	if nowBtc, err = m.BalanceOf(m.BaseCurrency); err != nil {
		log.Println(err.Error())
		return affected, err
	}

	//var needBtc float64
	for i := 0; i < len(m.Datas); i++ {
		if buyPrice < m.Datas[i].BuyPrice {
			buyIndex = i
			buyNeedBtc = m.Datas[buyIndex].TotalBaseCurrency
		} else {
			break
		}
	}
	//var needBtc float64
	for i := 0; i < len(m.Datas); i++ {
		if sellPrice < m.Datas[i].SellPrice {
			sellIndex = i
			sellNeedBtc = m.Datas[sellIndex].TotalBaseCurrency
		} else {
			break
		}
	}
	//补些btc
	if nowBtc < buyNeedBtc-(0.5*m.Datas[0].TotalBaseCurrency) && m.Datas[buyIndex].BuyOrder == "" {
		if err = m.BuyMarket((buyNeedBtc - nowBtc) * buyPrice); err != nil {
			log.Println(err.Error())
			return affected, err
		}
		if nowBtc, err = m.BalanceOf(m.BaseCurrency); err != nil {
			log.Println(err.Error())
			return affected, err
		}
	}
	//卖掉多出的btc
	if nowBtc > sellNeedBtc+(0.5*m.Datas[0].TotalBaseCurrency) && m.Datas[sellIndex].SellOrder == "" {
		if err = m.SellMarket(nowBtc - sellNeedBtc); err != nil {
			log.Println(err.Error())
			return affected, err
		}
		if nowBtc, err = m.BalanceOf(m.BaseCurrency); err != nil {
			log.Println(err.Error())
			return affected, err
		}
	}
	activeIndex := -1
	for i := 0; i < len(m.Datas); i++ {
		if nowBtc > m.Datas[i].TotalBaseCurrency-(0.5*m.Datas[0].TotalBaseCurrency) {
			activeIndex = i
		} else {
			break
		}
	}
	if activeIndex >= 0 {
		if m.Datas[activeIndex].SellOrder == "" {
			sellAmount := nowBtc
			if activeIndex > 0 { //四舍五入不精准,可能到时略有误差,最后一格检查是否有足够btc
				sellAmount = nowBtc - m.Datas[activeIndex-1].TotalBaseCurrency
			}
			orderid, err := m.SellOrder(utils.Digits(sellAmount, 6), m.Datas[activeIndex].SellPrice)
			if err != nil {
				log.Println(activeIndex, sellAmount, err.Error())
				return affected, err
			}
			m.Datas[activeIndex].SellOrder = orderid
			affected = true

			for j := 0; j < len(m.Datas); j++ {
				if j != activeIndex && m.Datas[j].SellOrder != "" {
					m.CancelOrder(m.Datas[j].SellOrder)
					m.Datas[j].SellOrder = ""
				}
			}
		}
	}
	if activeIndex < len(m.Datas)-1 {
		if m.Datas[activeIndex+1].BuyOrder == "" {
			buyAmount := m.Datas[activeIndex+1].TotalBaseCurrency - nowBtc
			orderid, err := m.BuyOrder(buyAmount, m.Datas[activeIndex+1].BuyPrice)
			if err != nil {
				log.Println(err.Error())
				return affected, err
			}
			m.Datas[activeIndex+1].BuyOrder = orderid
			affected = true
			//删除其他订单
			for j := 0; j < len(m.Datas); j++ {
				if j != activeIndex+1 && m.Datas[j].BuyOrder != "" {
					m.CancelOrder(m.Datas[j].BuyOrder)
					m.Datas[j].BuyOrder = ""
				}
			}
		}
	}
	return affected, nil
}
