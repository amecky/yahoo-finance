package yf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PriceInterval string

type DateRange string

const (
	PI_ONE_MINUTE     = PriceInterval("1m")
	PI_FIVE_MINUTES   = PriceInterval("5m")
	PI_THIRTY_MINUTES = PriceInterval("30m")
	PI_ONE_HOUR       = PriceInterval("1h")
	PI_FOUR_HOUR      = PriceInterval("4h")
	PI_ONE_DAY        = PriceInterval("1d")

	PR_ONE_DAY   = DateRange("1d")
	PR_ONE_MONTH = DateRange("1mo")
)

type YahooClient struct {
	Ticker    string
	URL       string
	DateRange DateRange
	Interval  PriceInterval
}

type YahooClientOption func(*YahooClient)

func WithSpecificDate(date string) YahooClientOption {
	return func(yc *YahooClient) {
		yc.updateUrl()
	}
}

func WithPriceInterval(pi PriceInterval) YahooClientOption {
	return func(yc *YahooClient) {
		yc.Interval = pi
		yc.updateUrl()
	}
}

func WithDateRange(dr DateRange) YahooClientOption {
	return func(yc *YahooClient) {
		yc.DateRange = dr
		yc.updateUrl()
	}
}

func (yc *YahooClient) updateUrl() {
	// period1=1659540600&period2=1660038300
	yc.URL = fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?symbol=ZAL.DE&range=%s&interval=%s&includePrePost=true&events=div%%7Csplit%%7Cearn&corsDomain=finance.yahoo.com", yc.Ticker, yc.DateRange, yc.Interval)
}
func NewYahooClient(ticker string, opts ...YahooClientOption) *YahooClient {
	ret := &YahooClient{
		Ticker:    ticker,
		Interval:  PI_ONE_DAY,
		DateRange: PR_ONE_MONTH,
		URL:       fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?range=%s&interval=%s&events=history&corsDomain=finance.yahoo.com", ticker, PR_ONE_MONTH, PI_ONE_DAY),
	}
	for _, o := range opts {
		o(ret)
	}
	return ret
}

type MetaData struct {
	Currency           string  `json:"currency"`
	Symbol             string  `json:"symbol"`
	ExchangeName       string  `json:"exchangeName"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
	DataGranularity    string  `json:"dataGranularity"`
	Range              string  `json:"range"`
}

type ChartData struct {
	Chart struct {
		Result []struct {
			Meta       MetaData `json:"meta"`
			Timestamps []int    `json:"timestamp"`
			Indicators struct {
				Quotes []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int     `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func (yc *YahooClient) Load() (MetaData, []Candle, error) {
	var md MetaData
	resp, err := http.Get(yc.URL)
	if err != nil {
		return md, nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return md, nil, errors.New(fmt.Sprintf("Statuscode: %d - Body: %s", resp.StatusCode, string(body)))
	}
	var chart ChartData
	json.Unmarshal(body, &chart)
	var ret = make([]Candle, 0)
	for _, r := range chart.Chart.Result {
		md = r.Meta
		for _, q := range r.Indicators.Quotes {
			for i := 0; i < len(q.Open); i++ {
				if q.Volume[i] > 0 {
					tm := time.Unix(int64(r.Timestamps[i]), 0)
					cnd := Candle{
						Timestamp: tm.Format("2006-01-02 15:04"),
						Open:      q.Open[i],
						High:      q.High[i],
						Low:       q.Low[i],
						Close:     q.Close[i],
						Volume:    q.Volume[i],
					}
					ret = append(ret, cnd)
				}
			}
		}
	}
	return md, ret, nil
}