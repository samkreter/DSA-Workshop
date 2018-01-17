from influxdb import InfluxDBClient
import time
from datetime import datetime, timezone

json_body = [
    {
        "measurement": "cpu_load_short",
        "tags": {
            "host": "server01",
            "region": "us-west"
        },
        "time": int(datetime.utcnow().timestamp()),
        "fields": {
            "value": 0.64
        }
    },
    {
        "measurement": "cpu_load_short",
        "tags": {
            "host": "server01",
            "region": "us-west"
        },
        "time": int(datetime.utcnow().timestamp()),
        "fields": {
            "value": 0.64
        }
    }
]

client = InfluxDBClient('localhost', 8086, 'root', 'root', 'example')

#client.create_database('example')

client.write_points(json_body)

result = client.query('select * from cpu_load_short;')

print("Result: {0}".format(result))
