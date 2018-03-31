package influxDB

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
)

type Tags map[string]string
type Fields map[string]interface{}

type InfluxdbClient struct {
	client client.Client
}

func New(username string, password string, host string) (*InfluxdbClient, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + host + ":8086",
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	influxdbClient := &InfluxdbClient{
		client: c,
	}

	return influxdbClient, nil
}

func (c *InfluxdbClient) Close() {
	c.client.Close()
}

func (c *InfluxdbClient) WritePoints(database string, name string, tags Tags, fields Fields, metricTime time.Time) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	point, err := client.NewPoint(
		name,
		tags,
		fields,
		metricTime,
	)
	if err != nil {
		return err
	}

	bp.AddPoint(point)

	err = c.client.Write(bp)
	if err != nil {
		return err
	}

	return nil
}

//fmt.Sprintf("select * from test_metric where location = '%s'", "westish")
func (c *InfluxdbClient) ReadMetrics(database string, query string) ([]models.Row, error) {
	q := client.Query{
		Command:  query,
		Database: database,
	}

	resp, err := c.client.Query(q)
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		log.Fatalln("Error: ", resp.Error())
	}

	return resp.Results[0].Series, nil
}
