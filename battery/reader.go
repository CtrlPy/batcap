package battery

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

type Reader interface {
	ReadInfo() (*BatteryInfo, error)
}
