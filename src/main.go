package main

import (
	"fmt"
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
}

