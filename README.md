# device-gps-go

## Introduction

This repository contains a GPS device service for EdgeX Foundry 1.0. It was built using the [EdgeX Device SDK Go](https://github.com/edgexfoundry/device-sdk-go) provided by the EdgeX community. 
Once up and running, the device service will automatically push the location, aka the  Geographical Coordinates from a **USB GPS receiver** to EdgeX, more specifically to the `core-data` microservice. This repository also contains a file with mock GPS data to allow users to test a GPS based application without recourse to a real GPS receiver. 

The GPS receiver is expected to support the [NMEA 0183 protocol](https://en.wikipedia.org/wiki/NMEA_0183), and more specifically in this instance using `BU-353-S4` which expects the GPS device to use a USB plug.  Additional information about NMEA, including how to decode the data, is available at http://aprs.gids.nl/nmea/. The current implementation extracts data from the `$GPRMC` prefixed sentence/line.

The data read from the GPS device is exported in JSON format and contains:

- Latitude
- Longitude
- Speed (in knots/hr)
- Unix timestamp as suppiled by the GPS device, and not the host machine (useful for time synchronization in the absence of network connectivity/[NTP](http://www.ntp.org/ntpfaq/NTP-s-def.htm)). 


## Getting Started

This tutorial assumes you have an EdgeX instance running on your Linux/Mac machine. This device service has not been tested on Windows yet. You will also need Go 1.11 installed.

### Device Service Using Mock GPS Data File

This repository comes with a file containing sample GPS output in `cmd/device-gps-go/gps_output_test.txt`. By Default, the GPS device profile will read data from this file.
When a real GPS device is connected, some connection specific detail needs to be provided to switch to real sensor mode.

1. Clone the repository in your go path. If using Go modules, install where you please.

```
$ cd ~/go/src/github.com/edgexfoundry/
$ git clone https://github.com/edgexfoundry-holding/device-gps
```

2. Build the service

```
$ cd device-gps-go
$ make build
```

3. Start the service

```
cd cmd/device-gps-go/
./device-gps-go
```

4. Once the device is running, go to your browser and enter the following URL to check if the device is generating data: `http://localhost:48080/api/v1/event/device/device-gps-go01/100`

### Device Service Using a Real BU-353-S4 Reciever

If you have a BU-353-S4 USB GPS Receiver, you can try reading the GPS data directly from the device. This installation requires some preparation.

1. Clone the repository in your go path. If using Go modules, install where you please.

```
$ cd ~/go/src/github.com/edgexfoundry/
$ git clone https://github.com/edgexfoundry-holding/device-gps
```

2. Plug your USB GSP device into your machine and check under `/dev/` which USB port corresponds to your device. For the sake of this tutorial, we will assume the GPS device is plugged into `/dev/ttyUSB0`.

3. The BU-353-S4 Device requires a 4800 baudrate. We need to assign that baud rate to the device.

```
$ stty 4800 > /dev/ttyUSB0
```

4. The default device service behavior is to read from the mock data file. This needs to be changed. Edit line 91 in `driver/simpledriver.go`, changing `gps_output_test.txt` to `/dev/ttyUSB0`.

5. Build the service. From the top level directory of the gps device service:

```
$ make build
```

6. If you previously ran the device service with the sample data, you need to remove the old device service from the database as well as all the associated readings.
 The easiest way to do this, if you don't have any other critical data in the database, is to simple clean it with a script provided by EdgeX. From anywhere on your machine (while EdgeX is running): 

```
$ wget https://raw.githubusercontent.com/edgexfoundry/developer-scripts/master/clean_mongo.js
$ mongo < clean_mongo.js
```

7. Start the service.

```
cd cmd/device-gps-go/
./device-gps-go
```

8. Once the device is running, go to your browser and enter the following URL to check if the device is generating data: `http://localhost:48080/api/v1/event/device/device-gps-go01/100`