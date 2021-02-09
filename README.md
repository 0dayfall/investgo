# investgo
investgo is a Golang package to retrieve data from Investing.com, which provides data retrieval from stocks funds, etfs currencies, indices, bonds, commodities, certificates and cryptocurrencies.

## Usage

It can be used from command line:
`go run cmd\main.go -symbolFile system.csv -fromDate "20150101" -toDate "20200209"` the results will be written to CSV files in /CSV, or a search can be made
`go run cmd\main.go -search APPL -country USA -assetType bond"`

or in code:
`investgo.GetSymbolHistoricalData(symbol, fromDate, toDate)` 
or like so:
`investgo.SearchSymbolJSON(symbol)`

## Thanks to
Thanks to investpy for inspiration https://pypi.org/project/investpy/
