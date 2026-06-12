package battery

import (
	"math"
	"testing"
	"time"
)

type mockReader struct {
	info *BatteryInfo
}

func (m *mockReader) ReadInfo() (*BatteryInfo, error) {
	return m.info, nil
}

func TestSessionIntegration(t *testing.T) {
	reader := &mockReader{
		info: &BatteryInfo{
			Status:     "Discharging",
			PowerNow:   10.0,
			EnergyNow:  50.0,
			Capacity:   50,
			CycleCount: 1,
		},
	}

	s := &Session{
		reader: reader,
		done:   make(chan struct{}),
		State: SessionState{
			LastUpdate:       time.Now(),
			LastPower:        10.0,
			IntegratedEnergy: 0,
		},
	}

	tickerChan := make(chan time.Time)
	s.ticker = &time.Ticker{C: tickerChan}

	// Use channel to synchronize instead of WaitGroup to prevent data races
	syncChan := make(chan struct{})
	s.SetUpdateCallback(func() {
		syncChan <- struct{}{}
	})

	go s.loop()

	// 1. Test stable 10W * 3600s (we will simulate 3600 ticks of 1 second)
	now := s.State.LastUpdate
	for i := 0; i < 3600; i++ {
		now = now.Add(1 * time.Second)
		tickerChan <- now
		<-syncChan // wait for loop to process
	}

	if math.Abs(s.State.IntegratedEnergy-10.0) > 0.01 {
		t.Errorf("Expected ~10 Wh, got %f Wh", s.State.IntegratedEnergy)
	}

	// 2. Test dt > 5.0 (skip integration)
	energyBefore := s.State.IntegratedEnergy
	now = now.Add(10 * time.Second) // huge gap
	tickerChan <- now
	<-syncChan

	if s.State.IntegratedEnergy != energyBefore {
		t.Errorf("Expected energy to not grow after sleep, but grew from %f to %f", energyBefore, s.State.IntegratedEnergy)
	}

	// 3. Status != "Discharging"
	reader.info.Status = "Charging"
	now = now.Add(1 * time.Second)
	tickerChan <- now
	<-syncChan

	if s.State.IntegratedEnergy != energyBefore {
		t.Errorf("Expected energy to not grow while charging, but grew")
	}

	s.Stop()
}
