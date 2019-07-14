package main

import "time"

type MarketPairDoc struct {
	Hour       time.Time             `bson:"hour"`
	MarketPair string                `bson:"market_pair"`
	Minutes    map[string]MarketData `bson:"minutes"`
}
