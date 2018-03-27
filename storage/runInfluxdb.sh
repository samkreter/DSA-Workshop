docker run \
    -p 8083:8083 \
    -p 8086:8086 \
    -v ~/dockerDir/influxdb-data:/var/lib/influxdb \
    influxdb