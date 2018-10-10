package ruuvi

import (
	"fmt"
	"sync"
	"encoding/binary"
	"log"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"

	"gitlab.com/jtaimisto/bluewalker/hci"
	"gitlab.com/jtaimisto/bluewalker/host"
	"gitlab.com/jtaimisto/bluewalker/ruuvi"
)

type Ruuvi struct {
	HciDevice     string
}

var sampleConfig = `
  ## Bluetooth device to use 
  # hci_device = "hci0"
`

// SampleConfig returns the documentation about the sample configuration
func (m *Ruuvi) SampleConfig() string {
	return sampleConfig
}

// Description returns a basic description for the plugin functions
func (m *Ruuvi) Description() string {
	return "Read nearby ruuvi tags"
}

func ruuviLoop(reportChan chan *host.ScanReport, acc telegraf.Accumulator) {
	for sr := range reportChan {
		for _, ads := range sr.Data {
			if ads.Typ == hci.AdManufacturerSpecific && 
					len(ads.Data) >= 2 && 
					binary.LittleEndian.Uint16(ads.Data) == 0x0499 {
				ruuviData, err := ruuvi.Unmarshall(ads.Data)
				if err != nil {
					log.Printf("Unable to parse ruuvi data: %s\n", err.Error())
				}
				tags := map[string]string{
					"device": sr.Address.String(),
				}

				fields := map[string]interface{}{
					"rssi": sr.Rssi,
					"humidity": ruuviData.Humidity,
					"temperature": ruuviData.Temperature,
					"pressure": ruuviData.Pressure,
					"voltage": ruuviData.Voltage,
					"acceleration_x": ruuviData.AccelerationX,
					"acceleration_y": ruuviData.AccelerationY,
					"acceleration_z": ruuviData.AccelerationZ,
				}
				acc.AddGauge("ruuvi", fields, tags, time.Now())
			}

		}
	}
}

// Gather is the main execution function for the plugin
func (m *Ruuvi) Gather(acc telegraf.Accumulator) error {
	if len(m.HciDevice) == 0 {
		return fmt.Errorf("Need hci device")
	}

	raw, err := hci.Raw(m.HciDevice)
	if err != nil {
		return fmt.Errorf("Error opening RAW HCI socket: %s - Run as root and disable bluetooth services", err.Error())
	}

	host := host.New(raw)
	if err = host.Init(); err != nil {
		host.Deinit();
		return fmt.Errorf("Unable to initialize host: %s", err.Error())
	}

	reportChan, err := host.StartScanning(false, nil);
	if err != nil {
		host.Deinit()
		return fmt.Errorf("Unable to start scanning: %s", err.Error());
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ruuviLoop(reportChan, acc)
		wg.Done()
	}()
	wg.Wait()
	return nil;
}


func init() {
	m := Ruuvi{}
	inputs.Add("ruuvi", func() telegraf.Input {
		m := m
		return &m
	})
}
