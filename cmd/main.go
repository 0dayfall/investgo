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
	var fromDate string
	var toDate string
	var symbolFile string
	flag.StringVar(&symbol, "symbol", "", "Symbol to download data for")
	flag.StringVar(&fromDate, "fromDate", "", "From date")
	flag.StringVar(&toDate, "toDate", "", "The date to obtain data to")
	flag.StringVar(&symbolFile, "symbolFile", "", "A file containing symbold to download data for")
	flag.Parse()

	if isFlagPassed("symbol") {

		investgo.GetSymbolHistoricalData(symbol, fromDate, toDate)
	}

	if isFlagPassed("symbolFile") {

		files, err := readLines(symbolFile)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file)
			investgo.SymbolHistoricalDataToCSV(file, fromDate, toDate)
		}
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
