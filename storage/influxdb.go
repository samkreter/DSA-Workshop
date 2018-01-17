package main

import (  
    "encoding/json"
    "fmt"
    "github.com/influxdata/influxdb/client/v2"
    "log"
    "math/rand"
    "time"
)

const (  
    database = "nodes"
    username = "root"
    password = "root"
)

var clusters = []string{"public", "private"}

func main() {  
    c := influxDBClient()
	writeTestMetrics(c)
	readTestMetrics(c)
}

func influxDBClient() client.Client {  
    c, err := client.NewHTTPClient(client.HTTPConfig{
        Addr:     "http://localhost:8086",
        Username: username,
        Password: password,
    })
    if err != nil {
        log.Fatalln("Error: ", err)
    }
    return c
}

func writeTestMetrics(c client.Client) {  
    bp, err := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  database,
        Precision: "s",
    })
    if err != nil {
        log.Fatal(err)
    }

	tags := map[string]string{
		"location": "westish",
	}

	fields := map[string]interface{}{
		"cpu_usage":  rand.Float64() * 100.0,
		"disk_usage": rand.Float64() * 100.0,
	}

	point, err := client.NewPoint(
		"test_metric",
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}

    bp.AddPoint(point)

    err = c.Write(bp)
    if err != nil {
        log.Fatal(err)
    }
}

func readTestMetrics(c client.Client) float64 {  
    q := client.Query{
        Command:  fmt.Sprintf("select * from test_metric where location = '%s'", "westish"),
        Database: database,
    }

    resp, err := c.Query(q)
    if err != nil {
        log.Fatalln("Error: ", err)
    }
    if resp.Error() != nil {
        log.Fatalln("Error: ", resp.Error())
    }

    res, err := resp.Results[0].Series[0].Values[0][1].(json.Number).Float64()
    if err != nil {
        log.Fatalln("Error: ", err)
    }

    return res
}