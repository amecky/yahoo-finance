package main

import (
	"fmt"

	"github.com/amecky/yahoo-finance/yf"
)

func main() {
	yc := yf.NewYahooClient("ZAL.DE", yf.WithPriceInterval(yf.PI_FIVE_MINUTES), yf.WithDateRange(yf.PR_ONE_DAY))
	md, candles, err := yc.Load()
	fmt.Println(md)
	if err == nil {
		for _, c := range candles {
			fmt.Println(c)
		}
	} else {
		fmt.Println(err)
	}
}
