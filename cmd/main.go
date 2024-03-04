package main

import (
	"fmt"
	"git.tcp.direct/kayos/sendkeys"
	u "github.com/bcicen/go-units"
	"github.com/hunterdyar/go-bluetooth-test/disto"
	"math"
	"tinygo.org/x/bluetooth"
)

func main() {
	d := disto.Disto{}
	callback := measure
	d.OnMeasure = &callback
	d.Connect(bluetooth.DefaultAdapter)
}

func measure(meters float64) {
	distance := u.NewValue(meters, u.Meter)
	feet := distance.MustConvert(u.Foot).Float()
	inches := feet * 12.0
	inches = math.Mod(inches, 12)
	feet = math.Floor(feet)
	fmt.Println(distance)
	fmt.Println("feet in:", feet, inches)
	output := fmt.Sprintf("%s", distance)
	k, err := sendkeys.NewKBWrapWithOptions(sendkeys.Noisy)
	if err != nil {
		println(err.Error())
		return
	}

	k.Type(output)
}
