package main

import (
	"fmt"
	"github.com/go-ble/ble"
	"github.com/ui-kreinhard/cleargrass-le/clearglass-le"
	"os"
	"os/signal"
)

func logHandleEnvironmentData(temperature clearglass_le.Temperature, humidity clearglass_le.Humidity, battery clearglass_le.Battery, addr ble.Addr) {
	fmt.Println("addr", addr, "temp", temperature, "humidity", humidity, "battery", battery)
}


func main() {
	//log.SetOutput(ioutil.Discard)

	clearGrass := clearglass_le.NewClearGreass(logHandleEnvironmentData)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		fmt.Println("CtrlC waiting")
		<-signalChan
		clearGrass.Stop()
	}()
	clearGrass.Init()
}
