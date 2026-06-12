package report

import (
	"fmt"
	"strings"
	"time"

	"batcap/battery"

	"github.com/charmbracelet/lipgloss"
)

func Generate(state battery.SessionState) string {
	bmsDiff := state.EnergyStart - state.EnergyCurrentBMS
	diffStr := ""
	diffVal := state.IntegratedEnergy - bmsDiff

	if bmsDiff > 0 {
		pct := (diffVal / bmsDiff) * 100.0
		if diffVal > 0 {
			diffStr = fmt.Sprintf("+%.2f Wh (+%.0f%%)", diffVal, pct)
		} else {
			diffStr = fmt.Sprintf("%.2f Wh (%.0f%%)", diffVal, pct)
		}
	} else {
		diffStr = fmt.Sprintf("%.2f Wh", diffVal)
	}

	avgPower := 0.0
	hours := state.ElapsedSeconds / 3600.0
	if hours > 0 {
		avgPower = state.IntegratedEnergy / hours
	}

	var sb strings.Builder

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2)

	bmsHealth := 0.0
	if state.EnergyFullDesign > 0 {
		bmsHealth = (state.EnergyFull / state.EnergyFullDesign) * 100.0
	}

	content := fmt.Sprintf(`      BATTERY CAPACITY REPORT      
───────────────────────────────────
 SYSTEM INFO
 Laptop:     %s
 Battery:    %s
 Cycles:     %d
 BMS Health: %.0f%% (%.2f Wh / %.2f Wh)
───────────────────────────────────
 Test duration:      %s
 Start charge:       %d%% (%.2f Wh)
 End charge:         %d%% (%.2f Wh)
                                  
 REAL capacity:      %.2f Wh      
 BMS reported:       %.2f Wh      
 Difference:         %s           
                                  
 Avg power draw:     %.1f W`,
		state.LaptopModel,
		state.BatteryModel,
		state.CycleCount,
		bmsHealth, state.EnergyFull, state.EnergyFullDesign,
		formatDuration(state.ElapsedSeconds),
		state.InitialCapacity, state.EnergyStart,
		state.CurrentCapacity, state.EnergyCurrentBMS,
		state.IntegratedEnergy,
		bmsDiff,
		diffStr,
		avgPower,
	)

	sb.WriteString(borderStyle.Render(content))

	return sb.String()
}

func formatDuration(sec float64) string {
	d := time.Duration(sec * float64(time.Second))
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	return fmt.Sprintf("%dh %dm", h, m)
}
