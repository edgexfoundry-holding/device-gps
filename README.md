# device-gps

## Introduction

This repository contains a GPS devices service for EdgeX Foundry 1.0. It was built using the `device-sdk-go` provided by the EdgeX community. 
Once up and running, the device service will automatically push Geographical Coordinates from a **USB GPS reciever** to EdgeX, more specifically to the `core-data` microservice. This repository also contains a file containing mock GPS to allow the user to test a GPS based application without purchasing a GPS reciever. 

This device profile requires a GPS reciever using the NMEA 0183 protocol. More specifically, I used a `BU-353-S4`, which speaks the NMEA 0183 protocol and has a USB plug. For more information on the NMEA 0183 protocol and how to decode it read this page: http://aprs.gids.nl/nmea/. The sentense this device reads is `$GPRMC`.