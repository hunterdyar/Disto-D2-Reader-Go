package main

import (
	"encoding/binary"
	"fmt"
	u "github.com/bcicen/go-units"
	"math"
	"strings"
	"tinygo.org/x/bluetooth"
)

var (
	adapter                  = bluetooth.DefaultAdapter
	measureDataUUID, _       = bluetooth.ParseUUID("3ab10100-f831-4395-b29d-570977d5bf94")
	measureCharacteristic, _ = bluetooth.ParseUUID("3ab10101-f831-4395-b29d-570977d5bf94")
)

func main() {
	println("enabling")

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

	var device *bluetooth.Device
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			println(err.Error())
			return
		}

		println("connected to ", result.Address.String())
	}

	// get services
	srvcs, err := device.DiscoverServices([]bluetooth.UUID{measureDataUUID})
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
	println("connected")

	char.EnableNotifications(func(buf []byte) {
		//todo: the conversion to int gives us more sigfig than given by the device.
		//this is because we had to go to float64 for the distance deal
		bits := binary.LittleEndian.Uint32(buf)
		meters32 := math.Float32frombits(bits)
		meters := float64(meters32)
		distance := u.NewValue(meters, u.Meter)
		feet := distance.MustConvert(u.Foot).Float()
		inches := feet * 12.0
		inches = math.Mod(inches, 12)
		meteFormat := u.FmtOptions{
			Label: true,
			Short: true,
			Precision: 3,
		}
		feet = math.Floor(feet)
		fmt.Printf("%s or %v feet %.3f inches\n",distance.Fmt(meteFormat),int(feet), inches)
		
	})

	select {}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

