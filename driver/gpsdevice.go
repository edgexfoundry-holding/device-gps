// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
// Copyright (C) 2019 VMware, Inc
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example implementation of
// a ProtocolDriver interface.
//
package driver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

type gpsData struct {
	Longitude float64
	Latitude  float64
	Time      int64
	Speed     float64
}

type GPSDevice struct {
	lc           logger.LoggingClient
	asyncCh      chan<- *dsModels.AsyncValues
	switchButton bool

	device  *os.File
	scanner *bufio.Scanner
	gpsdata string
}

// DisconnectDevice handles protocol-specific cleanup when a device
// is removed.
func (s *GPSDevice) DisconnectDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	return nil
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *GPSDevice) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	s.lc = lc
	s.asyncCh = asyncCh

	// Go routine reading the lines printed by the GPS device and saving the coordinates to s.gpsdata
	go func() {
		// Open device (or file with sample output)
		s.device, _ = os.Open("gps_output_test.txt")
		defer s.device.Close()

		s.scanner = bufio.NewScanner(s.device)

		// going through the lines in the scanner
		for s.scanner.Scan() {
			// checking if starts with "$GPRMC"
			splitLine := strings.Split(s.scanner.Text(), ",")
			if len(splitLine) > 0 {
				if splitLine[0] == "$GPRMC" {
					s.gpsdata = parseGPSline(splitLine)
					time.Sleep(1000 * time.Millisecond)
				}
			}
		}
	}()

	return nil
}

// Parse GPS string and extract: Longitude, latitude, speed and time.
// The data is stored in a GPSData struct and returned
// The *data* array contains the different data points read from the $GPRMC line
// data[1] => UTC Time. 220516 = 10:05:16 PM
// data[2] => Data status (A = OK, V = Warning)
// data[3] => Latitude
// data[4] => North or South
// data[5] => Longitude
// data[6] => East or West
// data[7] => Speed over ground in knots
// data[8] => Track made good in degrees True
// data[9] => UT date. 090419 = April 9th 2019
// data[10] => Magnetic variation degrees (Easterly var. subtracts from true course)
// data[11] => E or W variation
// data[12] => Checksum
func parseGPSline(data []string) string {

	// TODO: Handle errors
	timestamp, _ := time.Parse("020106150405", data[9]+data[1])

	latitude, _ := convertDegreesToDecimal(data[3], data[4])

	longitude, _ := convertDegreesToDecimal(data[5], data[6])

	speed, _ := strconv.ParseFloat(data[7], 64)

	// Create struct containing data extracted from the line
	gpsDataPoint := gpsData{
		Longitude: longitude,
		Latitude:  latitude,
		Time:      timestamp.Unix(),
		Speed:     speed,
	}

	// converting to json
	resp, err := json.Marshal(gpsDataPoint)

	if err != nil {
		fmt.Println(err)
	}

	return string(resp)
}

func convertDegreesToDecimal(degreesMinutes string, hemisphere string) (float64, error) {
	// 4916.45,N -> 49 deg, 16.45 minute north

	decimalIndex := strings.Index(degreesMinutes, ".")

	degrees, errDegrees := strconv.ParseFloat(degreesMinutes[:decimalIndex-2], 64)
	if errDegrees != nil {
		return 0.0, errDegrees
	}

	minutes, errMinutes := strconv.ParseFloat(degreesMinutes[decimalIndex-2:], 64)
	if errMinutes != nil {
		return 0.0, errMinutes
	}

	multiplier := 1.0
	if hemisphere == "S" || hemisphere == "W" {
		multiplier = -1.0
	}

	return multiplier * (degrees + (minutes / 60)), nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *GPSDevice) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	if len(reqs) != 1 {
		err = fmt.Errorf("GPSDevice.HandleReadCommands; too many command requests; only one supported")
		return
	}
	s.lc.Debug(fmt.Sprintf("GPSDevice.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))

	res = make([]*dsModels.CommandValue, 1)
	now := time.Now().UnixNano() / int64(time.Millisecond)
	if reqs[0].DeviceResourceName == "GPS" {
		cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, s.gpsdata)
		res[0] = cv
	}
	return
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *GPSDevice) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {
	err := fmt.Errorf("GPSDevice.HandleWriteCommands; this device does not support write commands.")
	return err
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *GPSDevice) Stop(force bool) error {
	s.lc.Debug(fmt.Sprintf("GPSDevice.Stop called: force=%v", force))
	return nil
}
