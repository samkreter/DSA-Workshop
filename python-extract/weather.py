import requests
import json
import logging
import time

API_URL = "https://api.coindesk.com/v1/bpi/currentprice.json"

def main():
    response = GetAPIResponse(API_URL)

    getDBObject(response)

    time.sleep(60)

def getDBObject(response):
    return {
        "code": response["bpi"]["USD"]["code"],
        "time": response["time"]["updatedISO"],
        "currRate": response["bpi"]["USD"]["rate_float"],
        "currRateString": response["bpi"]["USD"]["rate"]
    }

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


if __name__ == '__main__':
    main()
    

