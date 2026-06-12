package battery

import (
	"encoding/json"
	"os"
	"time"
)

type SessionState struct {
	StartTime         time.Time `json:"start_time"`
	LastUpdate        time.Time `json:"last_update"`
	EnergyStart       float64   `json:"energy_start"`
	EnergyCurrentBMS  float64   `json:"energy_current_bms"`
	IntegratedEnergy  float64   `json:"integrated_energy"` // How much energy actually discharged
	InitialCapacity   int       `json:"initial_capacity"`
	CurrentCapacity   int       `json:"current_capacity"`
	ElapsedSeconds    float64   `json:"elapsed_seconds"`
	PowerHistory      []float64 `json:"power_history"`
	LaptopModel       string    `json:"laptop_model"`
	BatteryModel      string    `json:"battery_model"`
	CycleCount        int       `json:"cycle_count"`
	EnergyFull        float64   `json:"energy_full"`
	EnergyFullDesign  float64   `json:"energy_full_design"`
}

type Session struct {
	reader Reader
	State  SessionState
	ticker *time.Ticker
	done   chan struct{}
	updateCb func()
	AutoStop chan struct{} // Channel to notify main app to stop
}

func NewSession(r Reader) (*Session, error) {
	s := &Session{
		reader:   r,
		done:     make(chan struct{}),
		AutoStop: make(chan struct{}),
	}
	return s, nil
}

func (s *Session) SetUpdateCallback(cb func()) {
	s.updateCb = cb
}

func (s *Session) StartOrResume(reset bool) error {
	if reset {
		os.Remove("/tmp/batcap-session.json")
	}

	if err := s.loadState(); err != nil {
		// New session
		info, err := s.reader.ReadInfo()
		if err != nil {
			return err
		}
		s.State = SessionState{
			StartTime:        time.Now(),
			LastUpdate:       time.Now(),
			EnergyStart:      info.EnergyNow,
			EnergyCurrentBMS: info.EnergyNow,
			IntegratedEnergy: 0,
			InitialCapacity:  info.Capacity,
			CurrentCapacity:  info.Capacity,
			PowerHistory:     make([]float64, 0),
			LaptopModel:      ReadLaptopModel(),
			BatteryModel:     info.Manufacturer + " " + info.ModelName,
			CycleCount:       info.CycleCount,
			EnergyFull:       info.EnergyFull,
			EnergyFullDesign: info.EnergyFullDesign,
		}
		s.saveState()
	}

	s.ticker = time.NewTicker(1 * time.Second)
	go s.loop()
	return nil
}

func (s *Session) loop() {
	for {
		select {
		case <-s.done:
			return
		case now := <-s.ticker.C:
			info, err := s.reader.ReadInfo()
			if err != nil {
				continue
			}
			
			dt := now.Sub(s.State.LastUpdate).Seconds()
			if dt < 0 {
				dt = 1.0 // safeguard
			}
			
			// Integrate: Power (W) * dt (seconds) / 3600 = Wh
			if info.Status == "Discharging" {
				energyWh := info.PowerNow * (dt / 3600.0)
				s.State.IntegratedEnergy += energyWh
			}

			s.State.LastUpdate = now
			s.State.EnergyCurrentBMS = info.EnergyNow
			s.State.CurrentCapacity = info.Capacity
			s.State.ElapsedSeconds += dt
			
			// Append power to history for graph (keep last 60 seconds)
			s.State.PowerHistory = append(s.State.PowerHistory, info.PowerNow)
			if len(s.State.PowerHistory) > 60 {
				s.State.PowerHistory = s.State.PowerHistory[1:]
			}

			s.saveState()

			if s.updateCb != nil {
				s.updateCb()
			}
			
			// Auto stop at 1% or less
			if info.Capacity <= 1 && info.Status != "Charging" {
				select {
				case s.AutoStop <- struct{}{}:
				default:
				}
			}
		}
	}
}

func (s *Session) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.done)
}

func (s *Session) saveState() {
	data, _ := json.Marshal(s.State)
	os.WriteFile("/tmp/batcap-session.json", data, 0644)
}

func (s *Session) loadState() error {
	data, err := os.ReadFile("/tmp/batcap-session.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.State)
}

func (s *Session) CurrentPower() float64 {
	info, err := s.reader.ReadInfo()
	if err != nil {
		return 0
	}
	return info.PowerNow
}
