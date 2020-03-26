package robot

import "log"

// work_sell 只卖手动买
func (m *GridBuy) work_buy() {
	var affacted bool = true
	var err error
	m.buy_init()
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
		if affacted, err = m.buy_clear_order(); affacted || err != nil {
			if err != nil {
				log.Println(err.Error())
			}
			continue
		}
		if affacted, err = m.buy_new_order(); affacted || err != nil {
			if err != nil {
				log.Println(err.Error())
			}
		}

	}
}

func (m *GridBuy) buy_init() error {
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
func (m *GridBuy) buy_clear_order() (affected bool, err error) {
	for i := 0; i < len(m.Datas); i++ {
		orderId := m.Datas[i].BuyOrder
		if orderId != "" { //买单结束取消
			orderReturn, err := m.GetOrder(orderId)
			if err != nil {
				return false, err
			}
			if orderReturn.Data.Createdat > 0 {
				m.Datas[i].SellOrder = ""
				affected = true
			}

		}

		orderId = m.Datas[i].SellOrder
		if orderId != "" { //卖单直接取消
			err = m.CancelOrder(orderId)
			if err != nil {
				return false, err
			}
			m.Datas[i].BuyOrder = ""
			affected = true
		}
	}
	return affected, nil
}

//只下买单
func (m *GridBuy) buy_new_order() (affected bool, err error) {
	//需要下买单的的最高高度
	buyIndex := 0
	//当前卖价
	nowPrice, err := m.BuyPrice()
	//当前btc额度
	nowBtc, err := m.BalanceOf(m.BaseCurrency)
	if err != nil {
		return affected, err
	}
	for i := 0; i < len(m.Datas); i++ {
		if nowPrice < m.Datas[i].BuyPrice {
			buyIndex = i + 1
		} else {
			break
		}
	}

	//补些btc
	if buyIndex > 0 && nowBtc < m.Datas[buyIndex-1].TotalBaseCurrency-(0.5*m.Datas[0].TotalBaseCurrency) {
		err = m.BuyMarket((m.Datas[buyIndex-1].TotalBaseCurrency - nowBtc) * nowPrice)
		affected = true
		return
	}

	//最后一格了,没得买了
	if buyIndex >= len(m.Datas) { //
		return affected, nil
	}
	for i := buyIndex; i < len(m.Datas); i++ {
		if nowBtc > m.Datas[i].TotalBaseCurrency-(0.5*m.Datas[0].TotalBaseCurrency) {
			continue
		}
		if m.Datas[i].BuyOrder == "" {
			buyAmount := m.Datas[i].TotalBaseCurrency - nowBtc
			orderid, err := m.BuyOrder(buyAmount, m.Datas[i].BuyPrice)
			if err != nil {
				return affected, err
			}
			m.Datas[i].BuyOrder = orderid
			affected = true
			//删除其他订单
			for j := 0; j < len(m.Datas); j++ {
				if j != i && m.Datas[j].BuyOrder != "" {
					m.CancelOrder(m.Datas[j].BuyOrder)
					m.Datas[j].BuyOrder = ""
				}
			}
		}
		break
	}
	return affected, nil
}

//手动市价卖btc
func (m *GridBuy) buy_action() (err error) {
	//需要下卖单的的最低高度
	sellIndex := -1
	//当前卖价
	sellPrice, err := m.SellPrice()
	//当前btc额度
	nowBtc, err := m.BalanceOf(m.BaseCurrency)
	if err != nil {
		return err
	}

	//<当前卖价所在的最低高度
	for i := 0; i < len(m.Datas); i++ {
		if sellPrice < m.Datas[i].SellPrice {
			sellIndex = i
		} else {
			break
		}
	}
	sellAmount := nowBtc
	if sellIndex >= 0 {
		sellAmount = nowBtc - m.Datas[sellIndex].TotalBaseCurrency
	}
	if sellAmount > m.Datas[0].TotalBaseCurrency*0.5 {
		err = m.SellMarket(nowBtc - m.Datas[sellIndex].TotalBaseCurrency)
	}
	return err
}
