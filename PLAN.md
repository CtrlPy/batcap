# batcap — Battery Capacity Measurement Utility

A TUI utility for Linux that measures the **real** battery capacity (Wh) by logging
discharge data from the BMS and comparing it to what the controller reports.

---

## Motivation

The BMS (Battery Management System) inside a laptop battery estimates capacity using
Coulomb counting and voltage tables — it never directly measures Wh via a controlled
discharge. `batcap` fills this gap by logging the discharge in real time and computing
the actual energy delivered.

---

## Tech Stack

- **Language:** Go
- **TUI framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Data source:** `/sys/class/power_supply/BAT0/` (no root required)

---

## Project Structure

```
batcap/
├── main.go
├── battery/
│   └── reader.go       # reads /sys/class/power_supply/
├── tui/
│   ├── model.go        # Bubble Tea model + update loop
│   ├── view.go         # TUI rendering
│   └── chart.go        # ASCII discharge graph
├── report/
│   └── report.go       # final report generation
└── go.mod
```

---

## How It Works

### Phase 1 — Startup
- Read initial BMS state: `energy_now`, `energy_full`, `voltage_now`, `current_now`
- Record `energy_start` as the baseline
- Warn the user if charge < 95% (measurement accuracy depends on starting near 100%)
- Warn if AC adapter is connected (discharge must happen on battery only)

### Phase 2 — Live TUI Monitoring
Refresh every **10 seconds**, display:
- Current charge: % and Wh in real time
- ASCII graph (X = elapsed time, Y = Wh remaining)
- Current power draw (W)
- Estimated time remaining

### Phase 3 — Final Report (triggered by Ctrl+C or battery cutoff)
Compute:
```
real_capacity = energy_start - energy_end
```
Display:

```
╔══════════════════════════════════╗
║     BATTERY CAPACITY REPORT      ║
╠══════════════════════════════════╣
║ Test duration:      2h 34m       ║
║ Start charge:       98% (50.4Wh) ║
║ End charge:         3%  (1.5Wh)  ║
║                                  ║
║ REAL capacity:      48.9 Wh      ║
║ BMS reported:       51.47 Wh     ║
║ Difference:        -2.57 Wh (-5%)║
║                                  ║
║ Avg power draw:     19.0 W       ║
╚══════════════════════════════════╝
```

---

## Measurement Accuracy Notes

- Accuracy is highest when starting from **100% charge** and running to **~3%**
- AC adapter **must be disconnected** during the entire test
- BMS typically has ±3–5% inherent error — this is expected
- The utility measures energy actually delivered, not theoretical cell capacity

---

## Battery Source Path

Default: `/sys/class/power_supply/BAT0`

Override via flag: `--battery BAT1`

---

## Files Read from sysfs

| File | Description |
|---|---|
| `energy_now` | Current charge in µWh |
| `energy_full` | BMS-estimated full capacity in µWh |
| `energy_full_design` | Factory design capacity in µWh |
| `voltage_now` | Current voltage in µV |
| `power_now` | Current power draw in µW |
| `capacity` | Percentage (0–100) |
| `status` | Charging / Discharging / Full |
| `cycle_count` | Total charge cycles |

---

## Target Platform

- Linux (tested on Fedora, should work on any distro)
- Requires read access to `/sys/class/power_supply/` (no root needed)
