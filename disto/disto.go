package disto

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"tinygo.org/x/bluetooth"
)

type Disto struct {
	Device    *bluetooth.Device
	Connected bool
	OnMeasure *func(float64)
}

var measureDataUUID, _ = bluetooth.ParseUUID("3ab10100-f831-4395-b29d-570977d5bf94")
var measureCharacteristic, _ = bluetooth.ParseUUID("3ab10101-f831-4395-b29d-570977d5bf94")

func (d *Disto) onReceiveData(buf []byte) {
	//todo: the conversion to int gives us more sigfig than given by the device.
	//this is because we had to go to float64 for the distance deal
	bits := binary.LittleEndian.Uint32(buf)
	meters32 := math.Float32frombits(bits)
	meters := float64(meters32)

	if d.OnMeasure != nil {
		f := *d.OnMeasure
		f(meters)
	} else {
		fmt.Println("Got Measure, no listeners")
	}
}

func (d *Disto) Connect(adapter *bluetooth.Adapter) {
	d.Connected = false
	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	ch := make(chan bluetooth.ScanResult, 1)

	// Start scanning.
	println("scanning...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		println("found:", result.Address.String(), result.LocalName())
		if strings.Contains(result.LocalName(), "DISTO") {
			adapter.StopScan()
			ch <- result
		}
	})

	select {
	case result := <-ch:
		d.Device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			println(err.Error())
			return
		}

		println("connected to ", result.Address.String())
	}

	// get services
	srvcs, err := d.Device.DiscoverServices([]bluetooth.UUID{measureDataUUID})
	must("discover services", err)

	if len(srvcs) == 0 {
		panic("could not find service. Is DISTO D2?")
	}

	service := srvcs[0]
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{measureCharacteristic})

	if err != nil {
		println(err)
	}

	if len(chars) == 0 {
		panic("could not find heart rate characteristic")
	}

	char := chars[0]
	d.Connected = true
	println("connected")

	char.EnableNotifications(d.onReceiveData)

	select {}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
