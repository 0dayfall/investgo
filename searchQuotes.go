package investgo

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func symbolId(symbol string, country string) int {
	s := SearchQuotes(symbol, country)

	for _, quote := range s.Quotes {
		if quote.Symbol == symbol && quote.Flag == country {
			return quote.PairId
		}
	}

	return 0
}

func SearchQuotes(symbol string, country string) Stock {
	data := url.Values{}
	data.Set("search_text", symbol)
	data.Set("tab", "quotes")
	data.Set("limit", "270")
	data.Set("offset", "0")

	investUrl := "https://www.investing.com"
	resource := "/search/service/SearchInnerPage"

	u, _ := url.ParseRequestURI(investUrl)
	u.Path = resource
	urlStr := u.String()

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	//request.PostForm = form

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	request.Header.Add("Accept", "text/html")
	request.Header.Add("Accept-Encoding", "gzip, deflate, br")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

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

	/*body, err := ioutil.ReadAll(reader)
	if err != nil && err != io.EOF {
		log.Fatalln(err)
	}
	fmt.Println(string(body))*/

	var s Stock
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(reader).Decode(&s)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

/*{"quotes":[{"pairId":1095913,"name":"Apple Inc","flag":"Italy","link":"\/equities\/apple-computer-inc?cid=1095913","symbol":"AAPLE","type":"Stock - Milan","pair_type_raw":"Equities","pair_type":"equities","countryID":10,"sector":0,"region":6,"industry":0,"isCrypto":false,"exchange":"Milan","exchangeID":6}],"total":{"quotes":1,"allResults":1},"filters":[]}*/
type Stock struct {
	Quotes []struct {
		PairId      int    `json:"pairId"`
		Name        string `json:"name"`
		Flag        string `json:"flag"`
		Link        string `json:"link"`
		Symbol      string `json:"symbol"`
		TypeString  string `json:"type"`
		PairType    string `json:"pair_type"`
		PairTypeRaw string `json:"pair_type_raw"`
		CountryID   int    `json:"countryID"`
		Sector      int    `json:"sector"`
		Region      int    `json:"region"`
		Industry    int    `json:"industry"`
		IsCrypto    bool   `json:"isCrypto"`
		Exchange    string `json:"exchange"`
		ExchangeID  int    `json:"exchangeID"`
	}
	Total struct {
		Quotes     int `json:"quotes"`
		AllResults int `json:"allResults"`
	}
	Filters []struct {
	}
}
