package influxDB

import (  
    "encoding/json"
    "github.com/influxdata/influxdb/client/v2"
    "log"
    "time"
)

type Tags map[string]string
type Fields map[string]interface{}

const (  
    database = "nodes"
    username = "root"
    password = "root"
)

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

func writePoints(c client.Client, name string, tags Tags, fields Fields, metricTime time.Time) {
    bp, err := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  database,
        Precision: "s",
    })
    if err != nil {
        log.Fatal(err)
    }

    point, err := client.NewPoint(
		name,
		tags,
		fields,
		metricTime,
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

//fmt.Sprintf("select * from test_metric where location = '%s'", "westish")
func readTestMetrics(c client.Client, query string) float64 {  
    q := client.Query{
        Command:  query,
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