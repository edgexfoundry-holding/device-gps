# device-gps-go

## Introduction

This repository contains a GPS devices service for EdgeX Foundry 1.0. It was built using the `device-sdk-go` provided by the EdgeX community. 
Once up and running, the device service will automatically push Geographical Coordinates from a **USB GPS reciever** to EdgeX, more specifically to the `core-data` microservice. This repository also contains a file containing mock GPS to allow the user to test a GPS based application without purchasing a GPS reciever. 

This device profile requires a GPS reciever using the NMEA 0183 protocol. More specifically, I used a `BU-353-S4`, which speaks the NMEA 0183 protocol and has a USB plug. For more information on the NMEA 0183 protocol and how to decode it read this page: http://aprs.gids.nl/nmea/. The sentense this device reads is `$GPRMC`.

The data is encoded in JSON format. Each data point contains a JSON object with the following information:

- Latitude
- Longitude
- Speed (in knots/hr)
- Unix timestamp as suppiled by the GPS device, and not the host machine (useful for keeping track of time without internet connection on resource constricted device)


## Getting Started

This tutorial assumes you have an EdgeX instance running on your Linux/Mac machine. This device service has not been tested on Windows yet. You will also need Go 1.11 installed.

### Device Service Using Mock GPS Data File

This repository comes with a file containing sample GPS output in `cmd/device-gps-go/gps_output_test.txt`. By Default, this device profile will read from this file.

1. Clone the repository in your go path.

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

If you have a BU-353-S4 USB GPS Reciever, you can try reading the GPS data directly from the device. This installation requires some preparation.

1. Clone the repository in your go path.

```
$ cd ~/go/src/github.com/edgexfoundry/
$ git clone https://github.com/edgexfoundry-holding/device-gps
```

2. Plug your USB GSP device into your machine and check under `/dev/` which USB port corresponds to your device. For the sake of this tutorial, we will assume the GPS device is plugged into `/dev/ttyUSB0`.

3. The BU-353-S4 Device requires a 4800 baudrate. We need to assign that baud rate to the device.

```
$ stty 4800 > /dev/ttyUSB0
```

4. The device service is currently settup to read the mock data file. You need to go change the file name from `gps_output_test.txt` to `/dev/ttyUSB0` on line 91 in `driver/simpledriver.go`

5. Build the service. From the top level directory of the gps device service:

```
$ make build
```

6. If you previoulsy ran the device service with the sample data, you need to remove the old device service from the database as well as all the associated readings. The easiest way to do this, if you don't have any other critical data in the database, is to simple clean it with a script provided by EdgeX. From anywhere on your machine (while EdgeX is running): 

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
