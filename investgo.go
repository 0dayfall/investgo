package investgo

import (
	"encoding/csv"
	"os"
)

func Search(symbol string) (Stock, error) {
	return searchQuotes(symbol)
}

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

func SearchJSon(symbol string, assetType string, country string) (string, error) {
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

func HistoricalDataToCSV(country string, assetType string, symbol string, fromDate string, toDate string) error {
	id, err := symbolId(country, assetType, symbol)
	if err != nil {
		return err
	}

	records, err := getStockHistoricalData(id, symbol, "01/01/2015", "02/01/2020", true, "ASC", "Daily")
	if err != nil {
		return err
	}

	return writeToCSV(symbol, records)

}

func GetHistoricalData(country string, assetType string, symbol string, fromDate string, toDate string) ([][]string, error) {
	var records [][]string

	id, err := symbolId(country, assetType, symbol)
	if err != nil {
		return records, err
	}

	records, err = getStockHistoricalData(id, symbol, "01/01/2015", "02/01/2020", true, "ASC", "Daily")
	if err != nil {
		return records, err
	}

	return records, nil
}

func writeToCSV(symbol string, records [][]string) error {

	if _, err := os.Stat("CSV"); os.IsNotExist(err) {
		mkdirErr := os.Mkdir("CSV", 0755)
		if mkdirErr != nil {
			return mkdirErr
		}
	} else {
		return err
	}

	//Write to csv file
	file, err := os.Create("CSV\\" + symbol + ".csv")
	if err != nil {
		return err
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)

	err = csvWriter.WriteAll(records)
	if err != nil {
		return err
	}

	csvWriter.Flush()

	return nil
}
