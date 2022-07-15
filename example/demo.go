package main

import (
	"dfTrade"
	"github.com/shopspring/decimal"
	"log"
)

var (
	userId   = "540950026469"
	password = "366300"
)

func main() {
	trade := dfTrade.NewDFTrade()

	for i := 0; i < 10; i++ {
		err := trade.Login(userId, password)
		if err != nil {
			log.Println(err)

		} else {
			break
		}
	}

	trade.QueryAssetAndPositionV1()
	trade.SendOrder("159607", "SA", decimal.NewFromFloat(0.897), 100, dfTrade.TradeTypeSell)
	trade.SendOrder("159607", "SA", decimal.NewFromFloat(0.897), 100, dfTrade.TradeTypeSell)
	trade.GetRevokeList()
	//trade.CancelOrder()
	//trade.QueryAssetAndPositionV1()
	//trade.QueryAssetAndPositionV1()
}
