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
POINT_BUFFER = 10
INFLUXDB_HOST = "localhost"
DATABASE_NAME = "mydb"
API_URL = "http://samkreter.com/metrics/{measurement}?start={start}&end={end}"

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


def main():
    #Connect to the Database
    client = InfluxDBClient(INFLUXDB_HOST, 8086, 'root', 'root', DATABASE_NAME)
    
    #Creates database if it doesn't exist
    client.create_database(DATABASE_NAME)

    measurement = "bitcoin"
    start_date = "2013-08-01"
    end_date = "2014-02-01"

    #Construct the api url
    url = API_URL.format(measurement = measurement, start = start_date, end = end_date)
    
    #Get the Json data from the API
    json_data = GetAPIResponse(url)
    if not json_data:
        sys.exit(1)

    logger.info("Processing and adding %d points to the databases.", len(json_data))

    buffered_count = 0
    point_buffer = []

    for data_point in json_data:
        #Tranform the JSON response into the correct format
        influx_point = CreateInfluxdbPoint(measurement, data_point['price'], data_point['timestamp'])
        
        point_buffer.append(influx_point)

        #If we have more points than the buffer, 
        # write all points to the DB and clear the buffer
        if len(point_buffer) > POINT_BUFFER:
            buffered_count += 1
            client.write_points(point_buffer)
            point_buffer = []

if __name__ == '__main__':
    main()

