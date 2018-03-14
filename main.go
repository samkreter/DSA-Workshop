package main

import (
	"log"

	"github.com/samkreter/DSA-Workshop/storage/influxdb"
	"strconv"
	"time"
	"bufio"
    "encoding/csv"
    "io"
    "os"
)

const (  
    database = "BitcoinPrice"
    username = "root"
    password = "root"
)

func main(){

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