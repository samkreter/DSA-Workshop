from influxdb import InfluxDBClient
import time
from datetime import datetime, timezone
import requests
import logging
import sys


INFLUXDB_HOST = "localhost"
DATABASE_NAME = "BitcoinPrice"


def main():
    #Connect to the Database
    client = InfluxDBClient(INFLUXDB_HOST, 8086, 'root', 'root', DATABASE_NAME)
    
    measurement = "Bitcoin"

    data_points = client.query("SELECT * FROM " + measurement)

    print(data_points)
    #TODO: Do something cool with the data

    print("DONE!")


if __name__ == '__main__':
    main()

