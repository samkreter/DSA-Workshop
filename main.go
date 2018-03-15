package main

import (
	"log"

	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samkreter/DSA-Workshop/storage/influxdb"
)

const (
	database            = "BitcoinPrice"
	username            = "root"
	password            = "root"
	BitcoinAPIUrl       = "https://api.coindesk.com/v1/bpi/historical/close.json?start=2013-09-01&end=2018-03-12"
	CurrencyAPITemplate = "http://data.fixer.io/api/{date}?access_key={api_key}&base={base}&symbols=usd"
	todaysDate          = "2018-03-13"
)

type BitcoinResp struct {
	Bpi DatePriceJson `json:"bpi"`
}

type CurrencyResp struct {
	Timestamp int64         `json:"timestamp"`
	Rates     DatePriceJson `json:"rates"`
}

type ArbJson map[string]interface{}

type DatePriceJson map[string]float64

func getCurrencyURL(apiKey string, date string, base string) string {
	r := strings.NewReplacer(
		"{api_key}", apiKey,
		"{date}", date,
		"{base}", base)

	return r.Replace(CurrencyAPITemplate)
}

type MetalIntervals struct {
	Intervals 	[]MetalInterval 	`json:"intervals"`
}

type MetalInterval struct {
	Timestamp	string 	`json:"start"`
	Price   	float64   `json:"open"`
}

func main(){
	c, err := influxDB.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	raw, err := ioutil.ReadFile("./python-extract/gold.json")
    if err != nil {
        log.Fatal(err)
    }

    var goldData MetalIntervals
    err = json.Unmarshal(raw, &goldData)
	if err != nil {
		log.Fatal(err)
	}

	for _, goldInterval := range goldData.Intervals{

		timestampInt, err := strconv.ParseInt(goldInterval.Timestamp, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		timestamp := time.Unix(timestampInt, 0)

		fields := influxDB.Fields{
			"price": goldInterval.Price,
		}

		err = c.WritePoints(database, "Gold", influxDB.Tags{}, fields, timestamp)
		if err != nil {
			log.Fatal(err)
		}

	}

	log.Println("Finished Loading into influxdb")
}


func ConvertDateToTime(date string) time.Time {
	layout := "2006-01-02"
	timestamp, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	return timestamp
}

func AddCurrencyApI() {
	c, err := influxDB.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	for i := 0; i < 497; i++ {
		date := ConvertDateToTime(todaysDate).AddDate(0, 0, 0-i).Format("2006-01-02")

		base := "MXN"

		url := getCurrencyURL("REDACTED", date, base)
		log.Fatal(url)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var jsonData CurrencyResp
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			log.Fatal(err)
		}

		for _, price := range jsonData.Rates {

			timestamp := time.Unix(jsonData.Timestamp, 0)

			fields := influxDB.Fields{
				"price": price,
			}

			err = c.WritePoints(database, base, influxDB.Tags{}, fields, timestamp)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Println("Finished Loading into influxdb")
}

func LoadHistoricalBitcoinFromAPI() {
	c, err := influxDB.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	resp, err := http.Get(BitcoinAPIUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData BitcoinResp
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	for datetime, price := range jsonData.Bpi {

		layout := "2006-01-02"
		timestamp, err := time.Parse(layout, datetime)
		if err != nil {
			panic(err)
		}

		fields := influxDB.Fields{
			"price": price,
		}

		err = c.WritePoints(database, "Bitcoin", influxDB.Tags{}, fields, timestamp)
		if err != nil {
			log.Fatal(err)
		}

	}

	log.Println("Finished Loading into influxdb")
}

func loadBitcoinCsvData() {

	c, err := influxDB.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	csvFile, err := os.Open("1coinUSD.csv")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		timestampInt, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		timestamp := time.Unix(timestampInt, 0)

		price, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		fields := influxDB.Fields{
			"price": price,
		}

		err = c.WritePoints(database, "Bitcoin", influxDB.Tags{}, fields, timestamp)
		if err != nil {
			log.Fatal(err)
		}

	}

	log.Println("Finished Loading into influxdb")
}
