// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018 IOTech Ltd
// Copyright (C) 2019 VMware, Inc
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a simple example of a device service.
package main

import (
	device_gps "github.com/edgexfoundry/device-gps"
	"github.com/edgexfoundry/device-gps/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-gps"
)

func main() {
	sd := driver.GPSDevice{}
	startup.Bootstrap(serviceName, device_gps.Version, &sd)
}
