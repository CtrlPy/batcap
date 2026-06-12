package battery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type BatteryInfo struct {
	EnergyNow        float64 // Wh
	EnergyFull       float64 // Wh
	EnergyFullDesign float64 // Wh
	PowerNow         float64 // W
	VoltageNow       float64 // V
	Capacity         int     // %
	Status           string  // Charging/Discharging/Full
	CycleCount       int
	ModelName        string
	Manufacturer     string
}

type Reader struct {
	path string
}

func NewReader(batteryName string) *Reader {
	return &Reader{
		path: filepath.Join("/sys/class/power_supply", batteryName),
	}
}

// ReadInt reads an integer value from a sysfs file.
func (r *Reader) ReadInt(filename string) (int64, error) {
	data, err := os.ReadFile(filepath.Join(r.path, filename))
	if err != nil {
		return 0, err
	}
	str := strings.TrimSpace(string(data))
	return strconv.ParseInt(str, 10, 64)
}

func (r *Reader) ReadString(filename string) (string, error) {
	data, err := os.ReadFile(filepath.Join(r.path, filename))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// ReadInfo reads all current stats and falls back to charge_* if energy_* is missing.
func (r *Reader) ReadInfo() (*BatteryInfo, error) {
	if _, err := os.Stat(r.path); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("battery path not found: %s", r.path)
	}

	info := &BatteryInfo{}
	
	// Read voltage first, needed for fallbacks
	voltageUv, _ := r.ReadInt("voltage_now")
	info.VoltageNow = float64(voltageUv) / 1_000_000.0

	// Helper to read energy directly or calculate from charge
	readEnergyOrCharge := func(energyName, chargeName string) float64 {
		if val, err := r.ReadInt(energyName); err == nil {
			return float64(val) / 1_000_000.0 // microWh to Wh
		}
		if val, err := r.ReadInt(chargeName); err == nil && info.VoltageNow > 0 {
			ah := float64(val) / 1_000_000.0
			return ah * info.VoltageNow // Wh
		}
		return 0
	}

	info.EnergyNow = readEnergyOrCharge("energy_now", "charge_now")
	info.EnergyFull = readEnergyOrCharge("energy_full", "charge_full")
	info.EnergyFullDesign = readEnergyOrCharge("energy_full_design", "charge_full_design")

	// Read power
	if val, err := r.ReadInt("power_now"); err == nil {
		info.PowerNow = float64(val) / 1_000_000.0
	} else if val, err := r.ReadInt("current_now"); err == nil && info.VoltageNow > 0 {
		amps := float64(val) / 1_000_000.0
		info.PowerNow = amps * info.VoltageNow
	}

	if cap, err := r.ReadInt("capacity"); err == nil {
		info.Capacity = int(cap)
	}

	if status, err := r.ReadString("status"); err == nil {
		info.Status = status
	}

	if count, err := r.ReadInt("cycle_count"); err == nil {
		info.CycleCount = int(count)
	}

	if model, err := r.ReadString("model_name"); err == nil {
		info.ModelName = model
	}

	if mfg, err := r.ReadString("manufacturer"); err == nil {
		info.Manufacturer = mfg
	}

	return info, nil
}

// ReadLaptopModel attempts to read the vendor and product name from DMI
func ReadLaptopModel() string {
	vendor, _ := os.ReadFile("/sys/devices/virtual/dmi/id/sys_vendor")
	product, _ := os.ReadFile("/sys/devices/virtual/dmi/id/product_name")
	
	vStr := strings.TrimSpace(string(vendor))
	pStr := strings.TrimSpace(string(product))
	
	if vStr != "" || pStr != "" {
		return strings.TrimSpace(vStr + " " + pStr)
	}
	return "Unknown Device"
}
