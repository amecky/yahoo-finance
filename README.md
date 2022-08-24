# yahoo-finance


Simple go library to download prices from yahoo.


## Client

## Options

### WithPriceInterval



### WithDateRange

### WithSpecificDate

### WithTimePeriod

## Example

The following is a code snippet to show how to download the prices for Tesla using 5 minutes interval
and for one day

```
yc := yf.NewYahooClient("TSLA", yf.WithPriceInterval(yf.PI_FIVE_MINUTES), yf.WithDateRange(yf.PR_ONE_DAY))
	md, candles, err := yc.Load()
	fmt.Println(md)
	if err == nil {
		for _, c := range candles {
			fmt.Println(c)
		}
	} else {
		fmt.Println(err)
	}
```
