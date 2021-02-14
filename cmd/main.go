package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	investgo "github.com/0dayfall/investgo"
)

func main() {

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	symbol := searchCmd.String("symbol", "", "Symbol name")
	country := searchCmd.String("country", "", "Country stock market")
	assetType := searchCmd.String("assetType", "", "The type of asset: equities, bond, etf, index, crypto")

	historicalCmd := flag.NewFlagSet("historical", flag.ExitOnError)
	historicalSymbol := historicalCmd.String("symbol", "", "Symbol name")
	historicalCountry := historicalCmd.String("country", "", "Country stock market")
	historicalAssetType := historicalCmd.String("assetType", "", "The type of asset: equities, bond, etf, index, crypto")
	historicalFromDate := historicalCmd.String("fromDate", "01/01/2015", "From date")
	historicalToDate := historicalCmd.String("toDate", time.Now().Format("02/01/2006"), "The date to obtain data to")

	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	fileName := fileCmd.String("symbols", "", "File name")
	fileCountry := fileCmd.String("country", "", "Country stock market")
	fileAssetType := fileCmd.String("assetType", "", "The type of asset: equities, bond, etf, index, crypto")
	fileFromDate := fileCmd.String("fromDate", "01/01/2015", "From date")
	fileToDate := fileCmd.String("toDate", time.Now().Format("02/01/2006"), "The date to obtain data to")

	switch os.Args[1] {
	case "search":
		searchCmd.Parse(os.Args[2:])

		jsonString, err := investgo.SearchJSON(*symbol, *assetType, *country)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(jsonString)

		os.Exit(0)

	case "historical":
		historicalCmd.Parse(os.Args[2:])

		records, err := investgo.GetHistoricalData(*historicalCountry, *historicalAssetType, *historicalSymbol, *historicalFromDate, *historicalToDate)
		if err != nil {
			log.Fatal(err)
		}

		printRecords(records)

		os.Exit(0)

	case "file":
		fileCmd.Parse(os.Args[2:])

		symbols, err := readLines(*fileName)

		if err != nil {
			log.Fatal(err)
		}

		for _, sym := range symbols {
			fmt.Printf("%s", sym)
			err := investgo.HistoricalDataToCSV(*fileCountry, *fileAssetType, sym, *fileFromDate, *fileToDate)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\tOK\n")
		}
		os.Exit(0)

	}

}

func printRecords(records [][]string) {
	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, set := range records {
		for _, record := range set {
			fmt.Fprintf(w, "%s\t", record)
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	/*
		for _, set := range records {
			for _, record := range set {
				fmt.Printf("\t%s", record)
			}
			fmt.Printf("\n")
		}
	*/
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

// isFlagPassed is used to check if a flag is set
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
