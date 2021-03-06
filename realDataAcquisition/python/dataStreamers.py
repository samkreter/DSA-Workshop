import requests
import logging
from datetime import datetime
from urllib.parse import urlencode

import apiHelper
import config


from influxdb import InfluxDBClient


BITCOIN_API_URL = "https://api.coindesk.com/v1/bpi/currentprice.json"
CURRENCY_API_URL = "https://api.fixer.io/latest?base=USD"
#For new token https://goldprice.com/

#https://api.coindesk.com/v1/bpi/historical/close.json?start=2013-09-01&end=2018-03-12

def putInDB(data):
    client = InfluxDBClient('localhost', 8086, 'root', 'root', 'example')

    client.write_points(data)

    for point in data:
        print(client.query('select * from ' + point["measurement"]))

def main():

    #Bitcoin
    #bitcoinData = apiHelper.GetAPIResponse(BITCOIN_API_URL, createBitcoinJSON)
    
    #Currencies
    #currencyData = apiHelper.GetAPIResponse(CURRENCY_API_URL, createCurrencyJSON)

    #Gold
    goldData = apiHelper.GetAPIResponse(getMetalURI(interval=1440, days=1000))
    print(goldData)

    #Silver
    #resp = apiHelper.GetAPIResponse(getMetalURI("Silver"))
    #silverData = createMetalJSON(resp, "Silver")

    #putInDB(currencyData)

def getMetalURI(metal="Gold", interval="1", days="1"):
    
    baseURI = "https://service.nfusionsolutions.biz/api/Metals/IntradaySpots?"

    params = {
        "token": config.GOLD_API_TOKEN,
        "metals": metal,
        "currency": "USD",
        "indicators": "%5B%5D",
        "interval": interval,
        "days": days
    }
    return baseURI + urlencode(params)

def createMetalJSON(response, metalType="Gold"):
    currInterval = response[0]['data']['intervals'][-1]

    #TIME: currInterval['start'][6:-7]

    return {
        "measurement": "metals",
        "tags": {
            "type": metalType,
        },
        "time": int(datetime.utcnow().timestamp()),
        "fields": {
            "value": currInterval['low']
        }
    }
    
def createCurrencyJSON(response):
    responses = []
    for currency, rate in response["rates"].items():
        responses.append(
        {
            "measurement": currency,
            "tags": {
                "base": response["base"],
            },
            "time": int(datetime.utcnow().timestamp()),
            "fields": {
                "value": rate
            }
        })
    return responses

def createBitcoinJSON(response):
    return {
        "measurement": "bitcoin",
        "tags": {
            "currentcy": response["bpi"]["USD"]["code"],
        },
        "time": response["time"]["updatedISO"],
        "fields": {
            "value": float(response["bpi"]["USD"]["rate_float"])
        }
    }




if __name__ == '__main__':
    main()
