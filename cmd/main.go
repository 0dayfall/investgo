package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	investgo "github.com/0dayfall/investgo"
)

func main() {
	var symbol string
	var jsonFormat bool
	var search string
	var fromDate string
	var toDate string
	var country string
	var assetType string
	var symbolFile string
	flag.StringVar(&symbol, "symbol", "", "Symbol to download data for")
	flag.BoolVar(&jsonFormat, "json", false, "Format as JSON")
	flag.StringVar(&search, "search", "", "Symbol to search for")
	flag.StringVar(&fromDate, "fromDate", "", "From date")
	flag.StringVar(&toDate, "toDate", "", "The date to obtain data to")
	flag.StringVar(&country, "country", "", "Country stock market")
	flag.StringVar(&assetType, "assetType", "", "The type of asset: equities, bond, etf, index, crypto")
	flag.StringVar(&symbolFile, "symbolFile", "", "A file containing symbold to download data for")
	flag.Parse()

	if isFlagPassed("symbolFile") {

		symbols, err := readLines(symbolFile)

		if err != nil {
			log.Fatal(err)
		}

		for _, symbol := range symbols {
			fmt.Println(symbol)
			err := investgo.HistoricalDataToCSV(country, assetType, symbol, fromDate, toDate)
			if err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}

	if isFlagPassed("search") && isFlagPassed("assetType") && isFlagPassed("country") && jsonFormat {
		jsonString, err := investgo.SearchJSon(search, assetType, country)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(jsonString)
		os.Exit(0)
	}

	if isFlagPassed("symbol") && isFlagPassed("assetType") && isFlagPassed("country") {

		records, err := investgo.GetHistoricalData(country, symbol, assetType, fromDate, toDate)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(records)
		os.Exit(0)
	}

	if isFlagPassed("search") && !jsonFormat {
		stocks, err := investgo.Search(search)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(stocks)
		os.Exit(0)
	}

	if isFlagPassed("search") && jsonFormat {
		jsonString, err := investgo.SearchSymbolJSON(search)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(jsonString)
		os.Exit(0)
	}

}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
