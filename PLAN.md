# batcap — Battery Capacity Measurement Utility

A TUI utility for Linux that measures the **real** battery capacity (Wh) by logging
discharge data from the BMS and integrating the consumed power over time.

---

## Motivation

The BMS (Battery Management System) inside a laptop battery estimates capacity using
Coulomb counting and voltage tables — it never directly measures Wh via a controlled
discharge. `batcap` fills this gap by logging the discharge in real time, integrating
the power consumption, and computing the actual energy delivered.

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
│   ├── reader.go       # reads /sys/class/power_supply/, handles energy vs charge
│   └── session.go      # timer logic (1s tick), power integration, persistence
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

### Phase 1 — Startup & Capabilities Detection
- Detect `/sys/class/power_supply/BAT0/`
- Determine if the battery reports `energy_*` (Wh) and `power_now` (W) or `charge_*` (Ah) and `current_now` (A). If the latter, it will automatically compute power using `voltage_now` (`W = A * V`).
- Read initial BMS state: `energy_full`, `voltage_now`, etc.
- Record `energy_start` as the BMS baseline.
- Warn the user if charge < 95% (measurement accuracy depends on starting near 100%).
- Warn if AC adapter is connected (discharge must happen on battery only).
- Resume from previous session if `/tmp/batcap-session.json` exists.

### Phase 2 — Live Measurement (Background)
Refresh every **1 second**, integrate:
- Read `power_now` (or compute from `current_now * voltage_now`).
- Integrate: `Total_Energy += Power * delta_t`.
- Save state to `/tmp/batcap-session.json` periodically to prevent data loss on suspend/crash.

### Phase 3 — Live TUI Monitoring
Refresh UI every **few seconds**, display:
- Current charge: % and Wh in real time (from BMS vs Integrated).
- ASCII graph (X = elapsed time, Y = Wh remaining).
- Current power draw (W).
- Estimated time remaining.

### Phase 4 — Final Report (triggered by Ctrl+C or battery cutoff)
Compute:
```
bms_diff_capacity = bms_energy_start - bms_energy_end
real_capacity = integrated_power_over_time
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

- Accuracy is highest when starting from **100% charge** and running to **~3%**.
- AC adapter **must be disconnected** during the entire test.
- BMS typically has ±3–5% inherent error — this is expected.
- The utility measures energy actually delivered using mathematical integration, bypassing BMS capacity estimation flaws.

---

## Battery Source Path

Default: `/sys/class/power_supply/BAT0`

Override via flag: `--battery BAT1`

---

## Files Read from sysfs

| File | Description | Fallback if missing |
|---|---|---|
| `energy_now` | Current charge in µWh | `charge_now` (µAh) |
| `energy_full` | BMS-estimated full capacity in µWh | `charge_full` (µAh) |
| `energy_full_design` | Factory design capacity in µWh | `charge_full_design` |
| `power_now` | Current power draw in µW | `current_now` (µA) |
| `voltage_now` | Current voltage in µV | Required for fallback |
| `capacity` | Percentage (0–100) | |
| `status` | Charging / Discharging / Full | |
| `cycle_count` | Total charge cycles | |

---

## Target Platform

- Linux (tested on Fedora, should work on any distro).
- Requires read access to `/sys/class/power_supply/` (no root needed).
