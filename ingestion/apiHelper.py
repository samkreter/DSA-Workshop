import requests
import logging


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

