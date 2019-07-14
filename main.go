package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPCallable interface {
	dial()
}

type HTTPCaller struct {
	endpoint      string
	client        *http.Client
	responsebytes *chan []byte
}

func (b HTTPCaller) dial() (err error) {
	glog.Info("Dialling...", b.endpoint)
	response, err := b.client.Get(b.endpoint)
	if err != nil {
		glog.Error(err)
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		glog.Error(err)
	}
	*b.responsebytes <- bodyBytes
	return
}

type MarketData struct {
	MarketName string
	High       float64
	Low        float64
	Volume     float64
	Created    time.Time
	Timestamp  time.Time
}

type HTTPConfig struct {
	HTTPTimeout time.Duration
}

type HTTPScraper struct {
	httpCallers      []HTTPCaller
	marketDataStream chan MarketData
	rawbytetream     chan []byte
	dbclient         *mongo.Client
}

type Scrapeable interface {
	Scrape()
}

func configure() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("yml")
	err := v.ReadInConfig() // Find and read the config file
	return v, err
}

func main() {
	v, err := configure()
	if err != nil {
		glog.Fatal(fmt.Errorf("error when reading config: %v", err))
	}
	httpScraper := NewHTTPScraper(v)
	httpScraper.Scrape()
}
