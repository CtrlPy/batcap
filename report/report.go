package report

import (
	"fmt"
	"strings"

	"github.com/CtrlPy/batcap/battery"

	"github.com/charmbracelet/lipgloss"
)

func Generate(state battery.SessionState) string {
	bmsDiff := state.EnergyStart - state.EnergyCurrentBMS
	diffStr := ""
	diffVal := state.IntegratedEnergy - bmsDiff

	if bmsDiff >= 0.1 {
		pct := (diffVal / bmsDiff) * 100.0
		if diffVal > 0 {
			diffStr = fmt.Sprintf("+%.2f Wh (+%.0f%%)", diffVal, pct)
		} else {
			diffStr = fmt.Sprintf("%.2f Wh (%.0f%%)", diffVal, pct)
		}
	} else {
		// If BMS diff is tiny or negative (charging), don't show confusing percentages
		if diffVal > 0 {
			diffStr = fmt.Sprintf("+%.2f Wh", diffVal)
		} else {
			diffStr = fmt.Sprintf("%.2f Wh", diffVal)
		}
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

	pctDropped := 0.0
	if state.EnergyFull > 0 {
		pctDropped = bmsDiff / state.EnergyFull
	}
	
	realHealthStr := "Not enough data (discharge at least 5%)"
	if pctDropped >= 0.05 && state.EnergyFullDesign > 0 {
		realFullCap := state.IntegratedEnergy / pctDropped
		realHealth := (realFullCap / state.EnergyFullDesign) * 100.0
		realHealthStr = fmt.Sprintf("%.0f%% (%.2f Wh / %.2f Wh)", realHealth, realFullCap, state.EnergyFullDesign)
	} else if pctDropped > 0 && state.EnergyFullDesign > 0 {
		realFullCap := state.IntegratedEnergy / pctDropped
		realHealth := (realFullCap / state.EnergyFullDesign) * 100.0
		realHealthStr = fmt.Sprintf("%.0f%% (Low accuracy, test too short)", realHealth)
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
                                  
 Avg power draw:     %.1f W
───────────────────────────────────
 CONCLUSION
 BMS Claimed Health: %.0f%%
 REAL TESTED HEALTH: %s`,
		state.LaptopModel,
		state.BatteryModel,
		state.CycleCount,
		bmsHealth, state.EnergyFull, state.EnergyFullDesign,
		battery.FormatDuration(state.ElapsedSeconds),
		state.InitialCapacity, state.EnergyStart,
		state.CurrentCapacity, state.EnergyCurrentBMS,
		state.IntegratedEnergy,
		bmsDiff,
		diffStr,
		avgPower,
		bmsHealth,
		realHealthStr,
	)

	sb.WriteString(borderStyle.Render(content))

	return sb.String()
}
