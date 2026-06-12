//go:build darwin
package battery

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
)

type DarwinReader struct {
}

func NewReader(batteryName string) Reader {
	return &DarwinReader{}
}

func (r *DarwinReader) ReadInfo() (*BatteryInfo, error) {
	cmd := exec.Command("ioreg", "-rn", "AppleSmartBattery")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	info := &BatteryInfo{
		Manufacturer: "Apple",
		ModelName:    "Internal Battery",
	}

	extractInt := func(key string) int64 {
		re := regexp.MustCompile(`"` + key + `"\s*=\s*(\d+)`)
		matches := re.FindSubmatch(out)
		if len(matches) > 1 {
			val, _ := strconv.ParseInt(string(matches[1]), 10, 64)
			return val
		}
		return 0
	}

	voltageMv := extractInt("Voltage") // mV
	amperageMa := extractInt("Amperage") // mA
    
	isCharging := bytes.Contains(out, []byte(`"IsCharging" = Yes`))
	isFull := bytes.Contains(out, []byte(`"FullyCharged" = Yes`))

	if isFull {
		info.Status = "Full"
	} else if isCharging {
		info.Status = "Charging"
	} else {
		info.Status = "Discharging"
	}

	amps := float64(amperageMa)
	if amps < 0 {
		amps = -amps
	}
	amps = amps / 1000.0

	volts := float64(voltageMv) / 1000.0
	info.VoltageNow = volts
	info.PowerNow = amps * volts

	currentMah := float64(extractInt("CurrentCapacity"))
	maxMah := float64(extractInt("MaxCapacity"))
	designMah := float64(extractInt("DesignCapacity"))

	info.EnergyNow = (currentMah / 1000.0) * volts
	info.EnergyFull = (maxMah / 1000.0) * volts
	info.EnergyFullDesign = (designMah / 1000.0) * volts

	if maxMah > 0 {
		info.Capacity = int((currentMah / maxMah) * 100.0)
	}

	info.CycleCount = int(extractInt("CycleCount"))

	return info, nil
}

func ReadLaptopModel() string {
	cmd := exec.Command("sysctl", "-n", "hw.model")
	out, err := cmd.Output()
	if err == nil {
		return string(bytes.TrimSpace(out))
	}
	return "Apple Mac"
}
