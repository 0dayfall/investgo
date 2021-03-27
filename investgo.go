package investgo

import (
	"encoding/csv"
	"os"
)

// Search is used to serch for a symbol
func Search(symbol string) (Stock, error) {
	return searchQuotes(symbol)
}

// SearchSymbolJSON is used to search
// and get the results in JSON format
func SearchSymbolJSON(symbol string) (string, error) {
	stock, err := searchQuotes(symbol)
	if err != nil {
		return "", err
	}

	jsonString, err := asJSON(stock)
	if err != nil {
		return "", err
	}

	return jsonString, nil
}

// SearchJSON searched with more parameters to eliminate to many hits
func SearchJSON(symbol string, assetType string, country string) (string, error) {
	stock, err := searchQuotesAssetTypeCountry(symbol, assetType, country)
	if err != nil {
		return "", err
	}

	jsonString, err := asJSON(stock)
	if err != nil {
		return "", err
	}

	return jsonString, nil
}

// HistoricalDataToCSV is used to write historical data to CSV
func HistoricalDataToCSV(country string, assetType string, symbol string, fromDate string, toDate string, dir string) error {
	id, err := symbolId(country, assetType, symbol)
	if err != nil {
		return err
	}

	records, err := getStockHistoricalData(id, symbol, fromDate, toDate, true, "ASC", "Daily")
	if err != nil {
		return err
	}

	return writeToCSV(symbol, records, dir)

}

// GetHistoricalData is used to get data in [][]string
// format as rows, colums: date, open, high, low, close, volume
func GetHistoricalData(country string, assetType string, symbol string, fromDate string, toDate string) ([][]string, error) {
	var records [][]string

	id, err := symbolId(country, assetType, symbol)

	if err != nil {
		return records, err
	}

	records, err = getStockHistoricalData(id, symbol, fromDate, toDate, true, "ASC", "Daily")
	if err != nil {
		return records, err
	}

	return records, nil
}

//Swapping Date, Close, Open, High, Low, Volume to
// Date, Open, High, Low, Close, Volume
func swapColumns(records [][]string) {

	for _, record := range records {
		close := record[1]
		open := record[2]
		high := record[3]
		low := record[4]

		//record[0] no swap
		record[1] = open
		record[2] = high
		record[3] = low
		record[4] = close
		//record[5] no swap

	}
}

func writeToCSV(symbol string, records [][]string, dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		mkdirErr := os.Mkdir(dir, 0755)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	//Write to csv file
	file, err := os.Create(dir + string(os.PathSeparator) + symbol + ".csv")
	if err != nil {
		return err
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)

	swapColumns(records)
	err = csvWriter.WriteAll(records)
	if err != nil {
		return err
	}

	csvWriter.Flush()

	return nil
}
