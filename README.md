[![Go Report Card](https://goreportcard.com/badge/github.com/steotia/go-analytics-crypto-scraper)](https://goreportcard.com/report/github.com/steotia/go-analytics-crypto-scraper) [![Maintainability](https://api.codeclimate.com/v1/badges/2735355a910e90726a53/maintainability)](https://codeclimate.com/github/steotia/go-analytics-crypto-scraper/maintainability)

# go-analytics-crypto-scraper

This is a library which gets cryptocurrency exchange rates for a pair of currencies. The data fetched is stored in MongoDB at a configurable 5 min interval. The calls are non blocking.

*NOTE*: Scraper starts AFTER 5 mins of server start or as configured in the `config.yml`

## Configuration

Check out the `config.yml` file for configuration options. Most notably,
```
scrape:
  marketpairs:
    - BTC-ADA
    - ETH-ADA
    - BTC-MUSIC
    - BTC-ETH
```
one can add or remove the currency pairs as required!

## Data Model
Each document obtained by hitting the summary API (Bittrex), does not create a unique document in Mongo. Instead, the time series data,
for every interval, is collected in an hourly document. So, essentially, if the period of scraping is 5 mins, 
then every hour has 1 document for an exchange pair rathen than have 12 documents, one for each period. The hourly document in MongoDB
looks like...
```
{
    "_id": ObjectID("5d2b7a5651aeeeffbb6feaf6"),
    "hour": ISODate("2019-07-14T18:00:00.000Z"),
    "market_pair": "ETH-ADA",
    "minutes": {
        "54": {
            "marketname": "ETH-ADA",
            "high": 0.00026171,
            "low": 0.00024418,
            "volume": 2469936.62183744,
            "created": ISODate("2017-11-28T17:28:32.077Z"),
            "timestamp": ISODate("2019-07-14T18:54:08.453Z")
        }
        ...
    }
}
```

You can visit http://localhost:8081/db/cryptocurrencies/marketvalues to have a look at persisted documents.

## Installation
The runtime is Dockerised, so the only requirement is having docker running on the system. So, just run

```docker-compose -f stack.yml up```

## Configuration
If there are any code or configuration changes, simply do 

```docker build -t go-analytics-crypto-scraper -f Dockerfile .``` 

to update the local docker image and then install as mentioned above.

## TODO
- Better code packaging! API code is better packaged.
- Scraper code quality can be better.
- Scraping starts N ticks AFTER server start, can change it to immediately start
- Tests!

