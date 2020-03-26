package robot

import "log"

// work_sell 只卖手动买
func (m *GridBuy) work_sell() {
	var affacted bool = true
	var err error
	m.sell_init()
	for {
		if affacted {
			if err = m.Save(); err != nil {
				log.Panicln(err.Error())
			}
		}
		if err = m.sleepOrStop(10); err != nil {
			m.Status = 0
			m.Istop = false
			return
		}
		if affacted, err = m.sell_clear_order(); affacted || err != nil {
			if err != nil {
				log.Println(err.Error())
			}
			continue
		}
		if affacted, err = m.sell_new_order(); affacted || err != nil {
			if err != nil {
				log.Println(err.Error())
			}
		}

	}
}

func (m *GridBuy) sell_init() error {
	//计算各个高度需要买的btc
	for i := 0; i < len(m.Datas); i++ {
		if i == 0 {
			m.Datas[i].TotalBaseCurrency = m.Datas[i].BuyMount / m.Datas[i].BuyPrice
		} else {
			m.Datas[i].TotalBaseCurrency = m.Datas[i-1].TotalBaseCurrency + m.Datas[i].BuyMount/m.Datas[i].BuyPrice
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
func (m *GridBuy) sell_clear_order() (affected bool, err error) {
	for i := 0; i < len(m.Datas); i++ {
		orderId := m.Datas[i].BuyOrder
		if orderId != "" { //买单直接取消
			err = m.CancelOrder(orderId)
			if err != nil {
				return false, err
			}
			m.Datas[i].BuyOrder = ""
			affected = true
		}

		orderId = m.Datas[i].SellOrder
		if orderId != "" { //卖单结算取消
			orderReturn, err := m.GetOrder(orderId)
			if err != nil {
				return false, err
			}
			if orderReturn.Data.Createdat > 0 {
				m.Datas[i].SellOrder = ""
				affected = true
			}
		}
	}
	return affected, nil
}

//只下卖单
func (m *GridBuy) sell_new_order() (affected bool, err error) {
	//需要下卖单的的最低高度
	sellIndex := -1
	//当前卖价
	nowPrice, err := m.SellPrice()
	//当前btc额度
	nowBtc, err := m.BalanceOf(m.BaseCurrency)
	if err != nil {
		return affected, err
	}
	//<当前卖价所在的最低高度
	for i := 0; i < len(m.Datas); i++ {
		if nowPrice < m.Datas[i].SellPrice {
			sellIndex = i
		} else {
			break
		}
	}
	if sellIndex == -1 { //卖完
		if nowBtc > m.Datas[0].TotalBaseCurrency*0.5 {
			m.SellMarket(nowBtc)
			affected = true
		}
		return affected, nil
	}
	//卖掉多出的
	if nowBtc > m.Datas[sellIndex].TotalBaseCurrency+m.Datas[0].TotalBaseCurrency*0.5 {
		m.SellMarket(nowBtc - m.Datas[sellIndex].TotalBaseCurrency)
		affected = true
		return affected, nil
	}
	//下卖单
	for i := sellIndex; i >= 0; i-- {
		if nowBtc < m.Datas[i].TotalBaseCurrency-m.Datas[0].TotalBaseCurrency*0.5 {
			continue
		}
		if m.Datas[i].SellOrder == "" {
			sellAmount := nowBtc

			if i > 0 { //四舍五入不精准,可能到时略有误差,最后一格检查是否有足够btc
				sellAmount = nowBtc - m.Datas[i-1].TotalBaseCurrency
			}
			if sellAmount < (m.Datas[0].TotalBaseCurrency)*0.5 { //判断账户已经没btc了
				break
			}
			orderid, err := m.SellOrder(sellAmount, m.Datas[i].SellPrice)
			if err != nil {
				return affected, err
			}
			m.Datas[i].SellOrder = orderid
			affected = true

			for j := 0; j < len(m.Datas); j++ {
				if j != i && m.Datas[j].SellOrder != "" {
					m.CancelOrder(m.Datas[j].SellOrder)
					m.Datas[j].SellOrder = ""
				}
			}

		}
		break
	}

	return affected, nil
}

//手动市价买btc
func (m *GridBuy) sell_action() (err error) {
	//需要下卖单的的最低高度
	buyIndex := -1
	//当前卖价
	nowPrice, err := m.BuyPrice()
	//当前btc额度
	nowBtc, err := m.BalanceOf(m.BaseCurrency)
	if err != nil {
		return err
	}
	//>当前买价所在的最低高度
	for i := 0; i < len(m.Datas); i++ {
		if nowPrice < m.Datas[i].BuyPrice {
			buyIndex = i
		} else {
			break
		}
	}
	if buyIndex == -1 {
		return
	}
	buyMount := m.Datas[buyIndex].TotalBaseCurrency - nowBtc
	if buyMount > m.Datas[0].TotalBaseCurrency*0.5 {
		err = m.BuyMarket(nowPrice * (m.Datas[buyIndex].TotalBaseCurrency - nowBtc))
	}
	return err
}
