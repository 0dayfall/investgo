package investgo

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

/*
	This function retrieves historical data from the introduced stock from Investing.com. So on, the historical data
    of the introduced stock from the specified country in the specified date range will be retrieved and returned as
    a :obj:`pandas.DataFrame` if the parameters are valid and the request to Investing.com succeeds. Note that additionally
    some optional parameters can be specified: as_json and order, which let the user decide if the data is going to
    be returned as a :obj:`json` or not, and if the historical data is going to be ordered ascending or descending (where the
    index is the date), respectively.
    Args:
        stock (:obj:`str`): symbol of the stock to retrieve historical data from.
        country (:obj:`str`): name of the country from where the stock is.
        from_date (:obj:`str`): date formatted as `dd/mm/yyyy`, since when data is going to be retrieved.
        to_date (:obj:`str`): date formatted as `dd/mm/yyyy`, until when data is going to be retrieved.
        as_json (:obj:`bool`, optional):
            to determine the format of the output data, either a :obj:`pandas.DataFrame` if False and a :obj:`json` if True.
        order (:obj:`str`, optional): to define the order of the retrieved data which can either be ascending or descending.
        interval (:obj:`str`, optional):
            value to define the historical data interval to retrieve, by default `Daily`, but it can also be `Weekly` or `Monthly`.
    Returns:
        :obj:`pandas.DataFrame` or :obj:`json`:
            The function can return either a :obj:`pandas.DataFrame` or a :obj:`json` object, containing the retrieved
            historical data of the specified stock from the specified country. So on, the resulting dataframe contains the
            open, high, low, close and volume values for the selected stock on market days and the currency in which those
            values are presented.
            The returned data is case we use default arguments will look like::
                Date || Open | High | Low | Close | Volume | Currency
                -----||------|------|-----|-------|--------|----------
                xxxx || xxxx | xxxx | xxx | xxxxx | xxxxxx | xxxxxxxx
            but if we define `as_json=True`, then the output will be::
                {
                    name: name,
                    historical: [
                        {
                            date: 'dd/mm/yyyy',
                            open: x,
                            high: x,
                            low: x,
                            close: x,
                            volume: x,
                            currency: x
                        },
                        ...
                    ]
                }
    Raises:
        ValueError: raised whenever any of the introduced arguments is not valid or errored.
        IOError: raised if stocks object/file was not found or unable to retrieve.
        RuntimeError: raised if the introduced stock/country was not found or did not match any of the existing ones.
        ConnectionError: raised if connection to Investing.com could not be established.
        IndexError: raised if stock historical data was unavailable or not found in Investing.com.
    Examples:
        >>> data = investpy.get_stock_historical_data(stock='bbva', country='spain', from_date='01/01/2010', to_date='01/01/2019')
        >>> data.head()
                     Open   High    Low  Close  Volume Currency
        Date
        2010-01-04  12.73  12.96  12.73  12.96       0      EUR
        2010-01-05  13.00  13.11  12.97  13.09       0      EUR
        2010-01-06  13.03  13.17  13.02  13.12       0      EUR
        2010-01-07  13.02  13.11  12.93  13.05       0      EUR
        2010-01-08  13.12  13.22  13.04  13.18       0      EUR
*/

func getStockHistoricalData(id int, stock string, fromDate string, toDate string, asJSON bool, order string, interval string) [][]string {

	/*requestBody, err := json.Marshal(map[string]string{
		"curr_id":      "482",
		"smlID":        "1159548",
		"header":       "ABB Historical Data",
		"st_date":      "01/01/2017",
		"end_date":     "01/01/2019",
		"interval_sec": "Daily",
		"sort_col":     "date",
		"sort_ord":     "DESC",
		"action":       "historical_data",
	})*/
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
		log.Fatalln(err)
	}
	//request.PostForm = form

	request.Header.Add("Accept", "text/plain, */*; q=0.01")
	request.Header.Add("Accept-Encoding", "gzip, deflate, br")
	//request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Referer", "https://www.investing.com/equities/abb-ltd-historical-data?cid=482")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	/*fmt.Println(request.Header)
	fmt.Println()
	fmt.Println(request.Body)
	fmt.Println()*/

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.Status)
	fmt.Println(resp.Header)

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
		log.Fatalln(err)
	}

	f, err := os.Create("logfile.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// attempt #1
	log.SetOutput(io.Writer(f))
	log.Println(string(body))

	// While have not hit the </html> tag

	//records := make([][]string, 0)
	records := [][]string{{"Date", "Open", "High", "Low", "Close", "Vol."}}
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

	return records

}
