package clearglass_le

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"log"
	"time"
)

type Temperature uint16 


func conv(data uint16) string {

	return fmt.Sprint(float32(uint16(data))/10.0)
}

func (t *Temperature) String() string {
	return conv(uint16(*t))
}

type Humidity uint16

func (h *Humidity) String() string {
	return conv(uint16(*h))
}

type Battery uint16
type OnTemperateChangeHandle func(temperature Temperature, humidity Humidity, battery Battery, addr ble.Addr)

type ClearGrass struct {
	handle OnTemperateChangeHandle
}

func (c *ClearGrass) onPeripheralDiscovered(a ble.Advertisement) {

	sds := a.ServiceData()
	for _, sd := range sds {
		//log.Print(a.Addr().String(), sd.UUID.String(), sd.Data)
		if len(sd.Data) > 9 && sd.UUID.String() == "fdcd" {
			serviceData := sd.Data
			temp := binary.LittleEndian.Uint16(serviceData[10:12])
			humidity := binary.LittleEndian.Uint16(serviceData[12:14])
			battery := sd.Data[16]
			c.handle(Temperature(temp), Humidity(humidity), Battery(battery), a.Addr())
		}
	}
}

func NewClearGreass(handler OnTemperateChangeHandle) *ClearGrass {
	ret := &ClearGrass{
		handler,
	}
	return ret
}

func (c *ClearGrass) OnTemperatureChange(handle OnTemperateChangeHandle) {
	c.handle = handle
}

func (c *ClearGrass) Init() error {
	log.Println("Scanning...")
	d, err := dev.NewDevice("default")
	if err != nil {
		return err
	}

	ble.SetDefaultDevice(d)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 99999*time.Hour))
	return ble.Scan(ctx, true, c.onPeripheralDiscovered, nil)
}

func (c *ClearGrass) Stop() error {
	log.Println("Stopping...")
	return ble.Stop()
}
