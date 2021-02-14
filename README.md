# investgo
investgo is a Golang package to retrieve data from Investing.com, which provides data retrieval from stocks funds, etfs currencies, indices, bonds, commodities, certificates and cryptocurrencies.

## Usage

### Command line
There are 3 command that can be run: search, historical and file. Search searches for bonds, equities, crypto and so on and presents results. The file command reads a file of symbold and write for each one historical data to SYMBOL.csv file. Here is an example of the historical command:

```
λ go run main.go historical -symbol ERICb -assetType equities -country Sweden -fromDate 02/03/2021
Date                    Open    High    Low     Close   Vol.
2021-02-12 01:00:00     113.60  113.50  113.75  112.25  4190852
2021-02-11 01:00:00     113.70  113.10  114.10  112.90  4408637
2021-02-10 01:00:00     113.05  113.55  113.95  112.35  7418479
2021-02-09 01:00:00     113.40  111.50  113.95  111.40  6876022
2021-02-08 01:00:00     111.50  111.05  111.85  109.80  5329644
2021-02-05 01:00:00     111.30  111.95  112.10  110.55  6902399
2021-02-04 01:00:00     111.85  110.20  111.95  110.20  9573125
2021-02-03 01:00:00     109.80  108.20  109.95  106.75  8809129
```

### Library
The golang functions can also be used as a library for another program:
```golang
records, err := investgo.GetHistoricalData(*historicalCountry, *historicalAssetType, *historicalSymbol, *historicalFromDate, *historicalToDate)
if err != nil {
	log.Fatal(err)
}

printRecords(records)
```

## Thanks to
Thanks to investpy for inspiration https://pypi.org/project/investpy/
