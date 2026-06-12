# batcap 🔋

`batcap` (Battery Capacity Measurement Utility) is a TUI tool for Linux that measures the **real, true capacity** of your laptop battery in Watt-hours (Wh) by continuously integrating the power consumption during a discharge cycle.

## Why batcap?
Most operating systems and built-in Battery Management Systems (BMS) estimate your battery health and remaining capacity using predefined voltage tables and basic Coulomb counting. Often, these estimates drift over time, leading to sudden drops in battery percentage (e.g., dropping from 20% to 0% instantly) or reporting a "Full Capacity" that isn't true.

`batcap` bypasses the BMS's internal capacity estimation. By polling the actual power draw (Watts) every second and mathematically integrating it over time, it calculates exactly how much energy (Watt-hours) the battery physically delivered to your laptop from 100% down to 0%.

## Features
- **Real Capacity Integration:** Computes `Energy (Wh) = ∫ Power(W) dt`, giving you the actual energy discharged.
- **Hardware Info & Fallbacks:** Automatically detects your laptop model, battery model, cycle count, and BMS Health. Works even if your Linux kernel only reports Charge (Ah) instead of Energy (Wh).
- **TUI Dashboard with Sparkline:** Clean, live-updating terminal interface that draws a live sparkline graph of your power draw using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Asciigraph](https://github.com/guptarohit/asciigraph).
- **Crash & Suspend Resilience:** Automatically saves state to `/tmp/batcap-session.json` every second. If your laptop unexpectedly shuts down, no data is lost!
- **Auto-Saving Reports:** Automatically saves a plain-text report to your working directory when the measurement is done, or if the battery hits 1% and the system is about to shut down.

## Installation
Make sure you have Go installed (1.20+ recommended).

```bash
git clone https://github.com/yourusername/batcap.git
cd batcap
go build -o batcap
sudo mv batcap /usr/local/bin/
```
*(Note: Root privileges are not required to run `batcap`, only read access to `/sys/class/power_supply` is needed, which is available to regular users).*

## Usage

1. Charge your laptop to **100%**.
2. Disconnect the AC power adapter.
3. Run `batcap` in your terminal:
   ```bash
   ./batcap
   ```
4. Leave the laptop running (you can continue using it, or leave it idle) until the battery is nearly empty (e.g., 3-5%).
5. Press `q` or `Ctrl+C` to stop the measurement and print the final report.

If you want to clear a previous interrupted session and start fresh:
```bash
./batcap --reset
```

If you have multiple batteries and want to measure a specific one:
```bash
./batcap --battery BAT1
```

## Example Report
```text
      BATTERY CAPACITY REPORT      
───────────────────────────────────
 SYSTEM INFO
 Laptop:     LENOVO 20QDCTO1WW
 Battery:    LGC 02DL004
 Cycles:     59
 BMS Health: 101% (51.47 Wh / 51.00 Wh)
───────────────────────────────────
 Test duration:      3h 45m
 Start charge:       100% (52.40 Wh)
 End charge:         5% (2.62 Wh)
                                  
 REAL capacity:      47.80 Wh      
 BMS reported:       49.78 Wh      
 Difference:         -1.98 Wh (-4%)           
                                  
 Avg power draw:     12.7 W
```

## How It Works
`batcap` reads data from `/sys/class/power_supply/BAT*/`. It prioritizes reading `power_now` and `energy_now`. If those are not exposed by your battery driver, it reads `current_now` and `voltage_now` to calculate the power draw dynamically. Every second, it adds `Power * (1/3600)` to its internal integrator.

## License
MIT
