package main

import (
	"log"

	"github.com/samkreter/DSA-Workshop/storage/influxdb"
	"strconv"
	"time"
	"bufio"
	"encoding/json"
	"encoding/csv"
    "io"
	"os"
	"net/http"
	"io/ioutil"
)

const (  
    database = "BitcoinPrice"
    username = "root"
    password = "root"
)

type BitcoinResp struct {
	Bpi 		BpiJson `json:"bpi"`
	// Disclaimer 	string 	`json:"disclaimer"`
	// Time 		ArbJson	`json:"time"`
}

type ArbJson map[string]interface{}

type BpiJson map[string]float64

func main(){
	c, err := influxDB.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	url := "https://api.coindesk.com/v1/bpi/historical/close.json?start=2013-09-01&end=2018-03-12"

	resp, err := http.Get(url)
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


func loadCsvData(){

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
		if err != nil{
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