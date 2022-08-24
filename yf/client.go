package yf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/amecky/fin-math/math"
)

type PriceInterval string

type DateRange string

type TimePeriodType int

const (
	PI_ONE_MINUTE     = PriceInterval("1m")
	PI_FIVE_MINUTES   = PriceInterval("5m")
	PI_THIRTY_MINUTES = PriceInterval("30m")
	PI_ONE_HOUR       = PriceInterval("1h")
	PI_FOUR_HOUR      = PriceInterval("4h")
	PI_ONE_DAY        = PriceInterval("1d")

	PR_ONE_DAY   = DateRange("1d")
	PR_ONE_WEEK  = DateRange("1wk")
	PR_ONE_MONTH = DateRange("1mo")
	PR_ONE_YEAR  = DateRange("1y")

	TPT_FIXED  = TimePeriodType(1)
	TPT_PERIOD = TimePeriodType(2)
	TPT_RANGE  = TimePeriodType(3)
)

type TimePeriod struct {
	Type   TimePeriodType
	First  string
	Second string
}

type YahooClient struct {
	Ticker     string
	URL        string
	TimePeriod TimePeriod
	Interval   PriceInterval
}

type YahooClientOption func(*YahooClient)

func WithSpecificDate(date string) YahooClientOption {
	return func(yc *YahooClient) {
		yc.TimePeriod = TimePeriod{
			Type:  TPT_FIXED,
			First: date,
		}
		yc.updateUrl()
	}
}

func WithPriceInterval(pi PriceInterval) YahooClientOption {
	return func(yc *YahooClient) {
		yc.Interval = pi
		yc.updateUrl()
	}
}

func WithTimePeriod(start, end string) YahooClientOption {
	return func(yc *YahooClient) {
		yc.TimePeriod = TimePeriod{
			Type:   TPT_PERIOD,
			First:  start,
			Second: end,
		}
		yc.updateUrl()
	}
}

func WithDateRange(dr DateRange) YahooClientOption {
	return func(yc *YahooClient) {
		yc.TimePeriod = TimePeriod{
			Type:  TPT_RANGE,
			First: "1wk",
		}
		yc.updateUrl()
	}
}

func convertToUnixTimestamp(date string) string {
	if date != "" {
		et, err := time.Parse("2006-01-02 15:04", date)
		if err == nil {
			return strconv.FormatInt(et.Unix(), 10)
		} else {

		}
	}
	return ""
}

func (yc *YahooClient) updateUrl() {
	// period1=1659540600&period2=1660038300
	t := yc.TimePeriod
	tp := "range=1d"
	if t.Type == TPT_PERIOD {
		tp = "range=" + yc.TimePeriod.First
	} else if t.Type == TPT_FIXED {
		st := convertToUnixTimestamp(t.First + " 00:00")
		et := convertToUnixTimestamp(t.First + " 23:59")
		tp = "period1=" + st + "&period2=" + et
	} else if t.Type == TPT_PERIOD {
		st := convertToUnixTimestamp(t.First + " 00:00")
		et := convertToUnixTimestamp(t.Second + " 23:59")
		tp = "period1=" + st + "&period2=" + et
	}
	yc.URL = fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?symbol=%s&%s&interval=%s&includePrePost=true&events=div%%7Csplit%%7Cearn&corsDomain=finance.yahoo.com", yc.Ticker, yc.Ticker, tp, yc.Interval)
}
func NewYahooClient(ticker string, opts ...YahooClientOption) *YahooClient {
	ret := &YahooClient{
		Ticker:   ticker,
		Interval: PI_ONE_DAY,
		TimePeriod: TimePeriod{
			Type:  TPT_FIXED,
			First: "1wk",
		},
		URL: fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?range=%s&interval=%s&events=history&corsDomain=finance.yahoo.com", ticker, PR_ONE_MONTH, PI_ONE_DAY),
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
	err = json.Unmarshal(body, &chart)
	if err != nil {
		return md, nil, err
	}
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

func (yc *YahooClient) LoadMatrix() (MetaData, *math.Matrix, error) {

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
	err = json.Unmarshal(body, &chart)
	if err != nil {
		return md, nil, err
	}
	ret := math.NewMatrix(6)
	for _, r := range chart.Chart.Result {
		md = r.Meta
		idx := 0
		for _, q := range r.Indicators.Quotes {
			for i := 0; i < len(q.Open); i++ {
				if q.Volume[i] > 0 {
					tm := time.Unix(int64(r.Timestamps[i]), 0)
					row := ret.AddRow(tm.Format("2006-01-02 15:04"))
					row.Set(math.OPEN, q.Open[idx])
					row.Set(math.HIGH, q.High[idx])
					row.Set(math.LOW, q.Low[idx])
					row.Set(math.CLOSE, q.Close[idx])
					row.Set(math.ADJ_CLOSE, q.Close[idx])
					row.Set(math.VOLUME, float64(q.Volume[idx]))
					idx++
				}
			}
		}
	}
	return md, ret, nil
}
