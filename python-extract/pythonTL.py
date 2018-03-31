from influxdb import InfluxDBClient
import time
from datetime import datetime, timezone
import requests
import logging

POINT_BUFFER = 10
API_URL = "samkreter.com/api/metrics/{measurement}?start={start}&end={end}"


def GetAPIResponse(url):
    try:
        response = requests.get(url)
    except:
        logging.error('Error while getting url.')
        return False

    if(response == None):
        logging.error('GetAPIResponse: No Request')
        return False

    if(response.status_code != 200):
        logging.error('GetAPIResponse: Status code not 200')
        return False

    return response.json()


def CreateInfluxdbPoint(measurement, price, timestamp):
    return {
            "measurement": measurement,
            "tags": {},
            "time": int(timestamp),
            "fields": {
                "price": price
            }
        }



def main():
    client = InfluxDBClient('localhost', 8086, 'root', 'root', 'mydb')

    measurement = "bitcoin"

    #Construct the api url
    url = API_URL.format(measurement = measurement, start = "2013-08-01", end = "2014-05-01")
    
    #Get the Json data from the API
    json_data = GetAPIResponse(url)

    logging.info("Processing and adding %d points to the databases.", len(json_data))

    point_buffer = []

    for data_point in json_data:
        #Tranform the JSON response into the correct format
        influx_point = CreateInfluxdbPoint(measurement, data_point['price'], data_point['timestamp'])
        
        point_buffer.append(influx_point)

        #If we have more points than the buffer, 
        # write all points to the DB and clear the buffer
        if len(point_buffer) > POINT_BUFFER:
            client.write_points(point_buffer)
            point_buffer = []

if __name__ == '__main__':
    main()

#client.create_database('example')
# result = client.query('select * from cpu_load_short;')
# print("Result: {0}".format(result))
