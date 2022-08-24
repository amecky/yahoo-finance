# yahoo-finance


Simple go library to download prices as candles from yahoo.


## Client

First you need to instantiate the client providing the ticker of the stock

```
yf.NewYahooClient(ticker string) *YahooClient
```
Then you can call the load method 
```
func (yc *YahooClient) Load() (MetaData, []Candle, error)
```
It will return some basic meta data and a slice of candles. In case of an error you will get the error of course.

## Using fin-math and the matrix

This is more an internal function. It builds a matrix from the github.com/amecky/fin-math project containing the data.

```
func (yc *YahooClient) LoadMatrix() (MetaData, *math.Matrix, error)
```
This matrix can be used to run all sorts of indicators and so on which are included in the fin-math repo.


## Candle

The client returns a slice of Candle. The Timestamp uses the format "yyyy-MM-dd HH:mm"

```
type Candle struct {
	Timestamp string
	High      float64
	Low       float64
	Open      float64
	Close     float64
	Volume    int
}
```

## Metadata

Just some basic informations

```
type MetaData struct {
	Currency           string
	Symbol             string
	ExchangeName       string
	RegularMarketPrice float64
	DataGranularity    string
	Range              string
}
```

## Options

You define use options as optional parameters to create the client. The following describes the possible options.

### WithPriceInterval

Use one of the defined intervals

* PI_ONE_MINUTE
* PI_FIVE_MINUTES
* PI_THIRTY_MINUTES
* PI_ONE_HOUR
* PI_FOUR_HOUR
* PI_ONE_DAY

Example

```
yf.NewYahooClient("TSLA", yf.WithPriceInterval(yf.PI_FIVE_MINUTES))
```

### WithDateRange

Use one of the defined date ranges

* PR_ONE_DAY
* PR_ONE_WEEK
* PR_ONE_MONTH
* PR_ONE_YEAR

Example
```
yf.NewYahooClient("TSLA", yf.WithDateRange(yf.PR_ONE_DAY))
```

### WithSpecificDate

Use this to get the data for a specific date in the format yyyy-MM-dd

Example
```
yf.NewYahooClient("TSLA", yf.WithSpecificDate("2022-08-24"))
```

### WithTimePeriod

Use this to get the data for a custom date range. Bothe dates have to be in the format yyyy-MM-dd

Example
```
yf.NewYahooClient("TSLA", yf.WithTimePeriod("2022-08-01","2022-08-24"))
```

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

## TODO

This list is meant for me

* add all supported intervals and ranges from yahoo
* include math.Matrix
* improve error handling when parsing json
* add more functionality from yahoo finance
