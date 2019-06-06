// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
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

type SimpleDriver struct {
	lc           logger.LoggingClient
	asyncCh      chan<- *dsModels.AsyncValues
	switchButton bool

	device  *os.File
	scanner *bufio.Scanner
	gpsdata string
}

// DisconnectDevice handles protocol-specific cleanup when a device
// is removed.
func (s *SimpleDriver) DisconnectDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	return nil
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *SimpleDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues) error {
	s.lc = lc
	s.asyncCh = asyncCh

	go func() {
		// Open device (or file with sample output)
		// Chane line 91 if reading from real device
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
func parseGPSline(data []string) string {

	timestamp, _ := time.Parse("020106150405", data[9]+data[1])

	latitude := convertDegreesToDecimal(data[3], data[4])

	longitude := convertDegreesToDecimal(data[5], data[6])

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

func convertDegreesToDecimal(degreesMinutes string, hemisphere string) float64 {
	// 4916.45,N -> 49 deg, 16.45 minute north

	decimalIndex := strings.Index(degreesMinutes, ".")

	degrees, _ := strconv.ParseFloat(degreesMinutes[:decimalIndex-2], 64)

	minutes, _ := strconv.ParseFloat(degreesMinutes[decimalIndex-2:], 64)

	if hemisphere == "N" || hemisphere == "E" {
		return degrees + (minutes / 60)
	}
	if hemisphere == "S" || hemisphere == "W" {
		return -(degrees + (minutes / 60))
	}

	return 0.0
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *SimpleDriver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	s.lc.Debug(fmt.Sprintf("test"))
	if len(reqs) != 1 {
		err = fmt.Errorf("SimpleDriver.HandleReadCommands; too many command requests; only one supported")
		return
	}
	s.lc.Debug(fmt.Sprintf("SimpleDriver.HandleReadCommands: protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))

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
func (s *SimpleDriver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	if len(reqs) != 1 {
		err := fmt.Errorf("SimpleDriver.HandleWriteCommands; too many command requests; only one supported")
		return err
	}
	if len(params) != 1 {
		err := fmt.Errorf("SimpleDriver.HandleWriteCommands; the number of parameter is not correct; only one supported")
		return err
	}

	s.lc.Debug(fmt.Sprintf("SimpleDriver.HandleWriteCommands: protocols: %v, resource: %v, parameters: %v", protocols, reqs[0].DeviceResourceName, params))
	var err error
	if s.switchButton, err = params[0].BoolValue(); err != nil {
		err := fmt.Errorf("SimpleDriver.HandleWriteCommands; the data type of parameter should be Boolean, parameter: %s", params[0].String())
		return err
	}

	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *SimpleDriver) Stop(force bool) error {
	s.lc.Debug(fmt.Sprintf("SimpleDriver.Stop called: force=%v", force))
	return nil
}
