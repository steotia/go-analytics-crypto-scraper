package main

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	httpTimeoutInSec = 10
	byteBufferSize   = 10
	jsonBufferSize   = 10
	scrapePeriod     = 1
)

func (e HTTPScraper) Scrape() {
	uptimeTicker := time.NewTicker(1 * time.Minute)
	for {
		select {
		// TODO Quit and close connections etc
		case <-uptimeTicker.C:
			for _, httpCaller := range e.httpCallers {
				go httpCaller.dial()
			}
		case marketData := <-e.marketDataStream:
			collection := e.dbclient.Database("cryptocurrencies").Collection("marketvalues")
			doc := MarketPairDoc{
				Hour:       marketData.Timestamp.Truncate(time.Hour),
				MarketPair: marketData.MarketName,
				Minutes: map[string]MarketData{
					strconv.Itoa((marketData.Timestamp.Minute())): marketData,
				},
			}
			filter := bson.M{
				"$and": bson.A{
					bson.M{"hour": bson.M{"$eq": doc.Hour}},
					bson.M{"market_pair": bson.M{"$eq": marketData.MarketName}},
				},
			}
			count, err := collection.CountDocuments(context.TODO(), filter)
			if err != nil {
				glog.Error(err)
			}
			glog.Info("Found:", count)
			updateResult, err := collection.UpdateOne(
				context.TODO(),
				filter,
				// bson.D{{"hour", doc.Hour}},
				bson.D{{"$set",
					bson.D{
						{
							"minutes." + strconv.Itoa((marketData.Timestamp.Minute())), marketData,
						},
					},
				}},
				options.Update().SetUpsert(false),
			)

			if err != nil {
				glog.Error(err)
			}
			glog.Info("Updated Document")
			if updateResult.MatchedCount == 0 {
				insertResult, err := collection.InsertOne(context.TODO(), doc)
				if err != nil {
					glog.Error(err)
				}
				glog.Info("Inserted a Single Document: ", insertResult.InsertedID)
				count, err := collection.CountDocuments(context.TODO(), filter)
				if err != nil {
					glog.Error(err)
				}
				glog.Info("Found:", count)
			}

		case bitrexBytes := <-e.rawbytetream:
			var bitrexJSON BitrexJSON
			err := json.Unmarshal([]byte(string(bitrexBytes)), &bitrexJSON)
			if err != nil {
				glog.Error(err)
				break
			}
			for _, result := range bitrexJSON.Result {
				timestamp, err := time.Parse("2006-01-02T15:04:05", result.TimeStamp)
				if err != nil {
					glog.Fatal(err)
				}
				created, err := time.Parse("2006-01-02T15:04:05", result.Created)
				if err != nil {
					glog.Fatal(err)
				}
				marketData := MarketData{
					MarketName: result.MarketName,
					High:       result.High,
					Low:        result.Low,
					Volume:     result.Volume,
					Timestamp:  timestamp,
					Created:    created,
				}
				e.marketDataStream <- marketData
			}
		}
	}
}

func setDefaults(config *viper.Viper) {
	config.SetDefault("scrape.timeout.http", httpTimeoutInSec)
	config.SetDefault("scrape.size.bytebuffer", byteBufferSize)
	config.SetDefault("scrape.size.jsonbuffer", jsonBufferSize)
	config.SetDefault("scrape.period", scrapePeriod)
}

func NewHTTPScraper(config *viper.Viper) *HTTPScraper {

	setDefaults(config)
	hTimeout := config.GetInt("scrape.timeout.http")
	bSize := config.GetInt("scrape.size.bytebuffer")
	jSize := config.GetInt("scrape.size.jsonbuffer")
	markets := config.GetStringSlice("scrape.marketpairs")
	scrapePeriod := config.GetInt("scrape.period")

	client, err := GetMongoDBClient()
	if err != nil {
		glog.Fatal(err)
	}

	bytetream := make(chan []byte, bSize)
	marketDataStream := make(chan MarketData, jSize)
	bitrexHTTPConfig := HTTPConfig{HTTPTimeout: time.Duration(hTimeout)}
	httpCallers := make([]HTTPCaller, len(markets))
	for i, m := range markets {
		httpCaller := NewBitrexHTTPCaller(m, bitrexHTTPConfig, &bytetream)
		httpCallers[i] = *httpCaller
	}

	return &HTTPScraper{
		httpCallers:      httpCallers,
		marketDataStream: marketDataStream,
		rawbytetream:     bytetream,
		dbclient:         client,
		scrapePeriod:     scrapePeriod,
	}
}
