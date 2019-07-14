# go-analytics-crypto-scraper

This is a library which gets cryptocurrency exchange rates for a pair of currencies. The data fetched is stored in MongoDB at a 5 min interval.

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
one can add or remove the currency pairs as required.

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

## Installation
The runtime is Dockerised, so the only requirement is having docker running on the system. 
