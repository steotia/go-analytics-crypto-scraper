package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const bitrexEndpoint = "https://api.bittrex.com/api/v1.1/public/getmarketsummary?market="

type BitrexResultJSON struct {
	MarketName     string
	High           float64
	Low            float64
	Volume         float64
	Last           float64
	BaseVolume     float64
	TimeStamp      string
	Bid            float64
	Ask            float64
	OpenBuyOrders  float64
	OpenSellOrders float64
	PrevDay        float64
	Created        string
}

type BitrexJSON struct {
	Success bool
	Message string
	Result  []BitrexResultJSON
}

func (bitrexJSON *BitrexJSON) parse(bodyBytes []byte) {
	json.Unmarshal([]byte(string(bodyBytes)), &bitrexJSON)
}

func NewBitrexHTTPCaller(market string, httpConfig HTTPConfig, byteChan *chan []byte) *HTTPCaller {
	var bitrexNetClient = &http.Client{
		Timeout: time.Second * httpConfig.HTTPTimeout,
	}
	endpoint := fmt.Sprintf("%s%s", bitrexEndpoint, market)
	bitrexHTTPCaller := HTTPCaller{
		endpoint:      endpoint,
		client:        bitrexNetClient,
		responsebytes: byteChan,
	}
	return &bitrexHTTPCaller
}
