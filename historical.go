package investgo

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

/* This function retrieves historical data from the introduced stock
   from Investing.com. So on, the historical data
   of the introduced stock from the specified country in the
   specified date range will be retrieved and returned as
        Date        Close  Open   High   Low    Volume
        2010-01-04  12.73  12.96  12.73  12.96       0
        2010-01-05  13.00  13.11  12.97  13.09       0
        2010-01-06  13.03  13.17  13.02  13.12       0
        2010-01-07  13.02  13.11  12.93  13.05       0
        2010-01-08  13.12  13.22  13.04  13.18       0
*/
func getStockHistoricalData(id int, stock string, fromDate string, toDate string, asJSON bool, order string, interval string) ([][]string, error) {

	//Watch out, it's COHL not OHLC
	records := [][]string{{"Date", "Close", "Open", "High", "Low", "Vol."}}

	form := url.Values{}
	form.Add("curr_id", strconv.Itoa(id))
	form.Add("smlID", "1159548")
	form.Add("header", stock+" Historical Data")
	form.Add("st_date", fromDate)
	form.Add("end_date", toDate)
	form.Add("interval_sec", interval)
	form.Add("sort_col", "date")
	form.Add("sort_ord", order)
	form.Add("action", "historical_data")

	url := "https://www.investing.com/instruments/HistoricalDataAjax"

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return records, err
	}

	request.Header.Add("Accept", "text/plain, */*; q=0.01")
	request.Header.Add("Accept-Encoding", "gzip, deflate, br")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Referer", "https://www.investing.com/equities/abb-ltd-historical-data?cid=482")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(request)

	if err != nil {
		return records, err
	}

	defer resp.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return records, err
	}

	row := make([]string, 6)
	k := 0
	tokenizer := html.NewTokenizer(strings.NewReader(string(body)))

	for {
		tokenizer.Next()
		t := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		if t.Data == "tbody" {
			break
		}
	}

	for loop := true; loop; {
		tt := tokenizer.Next()
		t := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {

		case html.StartTagToken:

			if t.Data == "td" {

				for _, attr := range t.Attr {
					if attr.Key == "data-real-value" {
						//fmt.Print()
						t := strings.TrimSpace(attr.Val)
						if k == 0 {
							i, err := strconv.ParseInt(t, 10, 64)
							if err != nil {
								panic(err)
							}
							tm := time.Unix(i, 0)
							t = tm.Format("2006-01-02 15:04:05")
						}
						row[k] = t
					}
				}

				k++
			}

		case html.EndTagToken:

			if t.Data == "tr" {
				records = append(records, row)
				row = make([]string, 6)
				k = 0
			}

			if t.Data == "table" {
				loop = false
				break
			}
		}
	}

	return records, nil
}
