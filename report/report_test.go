package report

import (
	"strings"
	"testing"

	"github.com/CtrlPy/batcap/battery"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name          string
		state         battery.SessionState
		expectedTexts []string
	}{
		{
			name: "Normal Case",
			state: battery.SessionState{
				EnergyStart:      50,
				EnergyCurrentBMS: 45, // bmsDiff = 5
				IntegratedEnergy: 5,
				EnergyFull:       50,
				EnergyFullDesign: 50,
			},
			expectedTexts: []string{
				"100% (50.00 Wh / 50.00 Wh)",
			},
		},
		{
			name: "Short Test (Low accuracy)",
			state: battery.SessionState{
				EnergyStart:      50,
				EnergyCurrentBMS: 49, // bmsDiff = 1, pctDropped = 0.02 (< 0.05)
				IntegratedEnergy: 1,
				EnergyFull:       50,
				EnergyFullDesign: 50,
			},
			expectedTexts: []string{
				"Low accuracy, test too short",
			},
		},
		{
			name: "Not enough data",
			state: battery.SessionState{
				EnergyStart:      50,
				EnergyCurrentBMS: 50, // bmsDiff = 0, pctDropped = 0
				IntegratedEnergy: 0,
				EnergyFull:       50,
				EnergyFullDesign: 50,
			},
			expectedTexts: []string{
				"Not enough data",
			},
		},
		{
			name: "EnergyFullDesign is 0",
			state: battery.SessionState{
				EnergyStart:      50,
				EnergyCurrentBMS: 40,
				IntegratedEnergy: 10,
				EnergyFull:       50,
				EnergyFullDesign: 0,
			},
			expectedTexts: []string{
				"Not enough data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Generate(tt.state)
			for _, text := range tt.expectedTexts {
				if !strings.Contains(result, text) {
					t.Errorf("Generate() did not contain %q", text)
				}
			}
		})
	}
}
