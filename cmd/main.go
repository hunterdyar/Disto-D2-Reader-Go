package main

import (
	"fmt"
	"math"

	"flag"
	"git.tcp.direct/kayos/sendkeys"
	u "github.com/bcicen/go-units"
	"github.com/hunterdyar/Disto-D2-Reader-Go/disto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"tinygo.org/x/bluetooth"
)

var logging = true
var keyboard = true

func main() {
	//flags setup
	flag.BoolVar(&logging, "l", false, "Log Measurements to JSON output")
	flag.BoolVar(&keyboard, "k", false, "Type measurement as keyboard. May require raised permissions.")

	flag.Parse()
	d := disto.Disto{}
	callback := measure
	d.OnMeasure = &callback
	d.Connect(bluetooth.DefaultAdapter)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func measure(meters float64) {
	distance := u.NewValue(meters, u.Meter)
	feet := distance.MustConvert(u.Foot).Float()
	inches := feet * 12.0
	inches = math.Mod(inches, 12)
	feet = math.Floor(feet)

	output := fmt.Sprintf("%s", distance)

	if logging {
		log.Info().Float64("meters", meters).Msg("")

	}
	if keyboard {
		k, err := sendkeys.NewKBWrapWithOptions(sendkeys.Noisy)
		if err != nil {
			println(err.Error())
			return
		}
		k.Type(output)
	}

	if !keyboard && !logging {
		fmt.Printf("%f\n", meters)
	}
}
