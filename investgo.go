package investgo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func SymbolHistoricalDataToCSV(symbol string, fromDate string, toDate string) {
	id := symbolId(symbol, "Sweden")

	records := getStockHistoricalData(id, symbol, "01/01/2015", "02/01/2020", true, "ASC", "Daily")

	writeToCSV(symbol, records)
}

func GetSymbolHistoricalData(symbol string, fromDate string, toDate string) [][]string {
	id := symbolId(symbol, "Sweden")
	records := getStockHistoricalData(id, symbol, "01/01/2015", "02/01/2020", true, "ASC", "Daily")

	return records
}

func writeToCSV(symbol string, records [][]string) {
	fmt.Println(symbol)

	if _, err := os.Stat("CSV"); os.IsNotExist(err) {
		os.Mkdir("CSV", 0755)
	}

	//Write to csv file
	file, err := os.Create("CSV\\" + symbol + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	err = csvWriter.WriteAll(records)

	if err != nil {
		log.Fatal(err)
	}

	csvWriter.Flush()
}
