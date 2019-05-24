# device-gps

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

### Device Service Using Mock GPS DATA

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

