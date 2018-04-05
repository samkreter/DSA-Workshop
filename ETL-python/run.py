from influxdb import InfluxDBClient
import time
from datetime import datetime, timezone
import requests
import logging
import sys

#Set up logging stuff
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

#Base defaults. No worries on changing these
POINT_BUFFER = 1000 #Number of points to buffer before writtening them to INfluxdb
INFLUXDB_HOST = "localhost"
DATABASE_NAME = "BitcoinTest"
API_URL = "http://samkreter.com/api/metrics/{measurement}?start={start}&end={end}"

def GetAPIResponse(url):
    try:
        response = requests.get(url)
    except Exception as e:
        logger.error('Failed to get from URL: ' + str(e))
        return False

    if(response == None):
        logger.error('GetAPIResponse: No Request')
        return False

    if(response.status_code != 200):
        logger.error('GetAPIResponse: Non 200 status code of ' +
                     str(response.status_code))
        return False

    return response.json()


def CreateInfluxdbPoint(measurement, price, timestamp):
    return {
            "measurement": measurement,
            "tags": {},
            "time": timestamp,
            "fields": {
                "price": price
            }
        }


def AddDataToInfluxdb(client, json_data, measurement):
    buffered_count = 0
    point_buffer = []

    logger.info("Processing and adding %d points to the databases.", len(json_data))

    for data_point in json_data:
        #Tranform the JSON response into the correct format
        influx_point = CreateInfluxdbPoint(measurement, float(data_point['price']), data_point['timestamp'])
        
        point_buffer.append(influx_point)

        #If we have more points than the buffer, 
        # write all points to the DB and clear the buffer
        if len(point_buffer) > POINT_BUFFER:
            print("Adding buffered Points: ", buffered_count)
            buffered_count += 1
            client.write_points(point_buffer)
            point_buffer = []
    
    # Make sure to add any points left over that were less than the buffer
    if  len(point_buffer) > 0:
        client.write_points(point_buffer)


def main():
    #Connect to the Database
    client = InfluxDBClient(INFLUXDB_HOST, 8086, 'root', 'root', DATABASE_NAME)
    
    #Creates database if it doesn't exist
    client.create_database(DATABASE_NAME)

    measurement = "bitcoin"

    dates = ["2013-01-01", "2014-01-01", "2015-01-01", "2016-01-01", "2017-01-01", "2018-01-01"]

    json_data = []

    #Loop through the list of start and end dates
    for i in range(len(dates) - 1):
        
        logger.info("Getting data for start: " + dates[i] + " end: " + dates[i + 1])
        
        #Construct the api url
        url = API_URL.format(measurement = measurement, start = dates[i], end = dates[i + 1])
        
        #Get the Json data from the API
        json_data = GetAPIResponse(url)
        if not json_data:
            sys.exit(1)

        AddDataToInfluxdb(client, json_data, measurement)

if __name__ == '__main__':
    main()

