# batcap

`batcap` is a small terminal tool for checking the practical battery capacity of a laptop.

I wrote it mostly for myself, because I wanted to better understand the real condition of the battery in my old laptop. The system can show battery health and full charge capacity, but I wanted to run a simple full-discharge test and see how much energy the battery actually delivers in practice.

The tool watches battery power usage over time and calculates the discharged energy in watt-hours (Wh).

It is not a laboratory-grade measurement tool.
It is a small practical utility for people who like to test, compare, and better understand their laptop batteries.

## Why I made it

Old laptop batteries can be confusing.

Sometimes the system says the battery is still healthy, but the laptop suddenly drops from 20% to 0%. Sometimes the reported full capacity looks fine, but the real runtime feels worse than expected.

I wanted a simple way to answer one question:

**How much usable energy can this battery actually deliver during a discharge test?**

So `batcap` measures the battery while the laptop is running on battery power and produces a small report at the end.

## What it does

* Tracks battery power draw while the laptop is discharging.
* Calculates discharged energy in Wh.
* Shows a live terminal dashboard.
* Saves the test session while running.
* Generates a plain text report.
* Shows basic battery info such as model, cycle count, and reported health when available.
* Can reset an interrupted session and start fresh.
* Can select a specific battery if the laptop has more than one.

## Installation

### go install

If you have Go installed:

```bash
go install github.com/CtrlPy/batcap@latest
```

The binary will be placed in your `$GOPATH/bin` (or `$HOME/go/bin` by default).

### Download a release

Download the latest binary from the [Releases](https://github.com/CtrlPy/batcap/releases) page:

```bash
tar -xzf batcap_Linux_x86_64.tar.gz
sudo mv batcap /usr/local/bin/
```

### Build from source

```bash
git clone https://github.com/CtrlPy/batcap.git
cd batcap
go build -o batcap
sudo mv batcap /usr/local/bin/
```

On Linux, `batcap` reads battery data from `/sys/class/power_supply`, so it usually does not need root permissions to run.

## Usage

Charge your laptop to 100%, unplug the charger, and run:

```bash
batcap
```

Then let the laptop discharge.

You can use the laptop normally, but for cleaner comparison between tests, it is better to keep the workload similar each time.

When you want to stop the test, press:

```text
q
```

or:

```text
Ctrl+C
```

To clear an old interrupted session and start fresh:

```bash
batcap --reset
```

To test a specific battery:

```bash
batcap --battery BAT1
```

## Example report

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
 Measured capacity:  47.80 Wh
 BMS reported:       49.78 Wh
 Difference:         -1.98 Wh (-4%)

 Avg power draw:     12.7 W
───────────────────────────────────
 CONCLUSION
 Reported health:    100%
 Tested health:      96% (47.80 Wh / 50.00 Wh)
```

## How it works

`batcap` periodically reads battery data from the system and tracks power usage during the discharge test.

In simple words:

```text
energy used = power draw × time
```

The tool sums this over the whole test and shows the result in Wh.

On Linux, it reads data from:

```text
/sys/class/power_supply/BAT*/
```

On macOS, it uses system battery information from `ioreg`.

## A few notes

The result depends on the test conditions.

Screen brightness, CPU load, Wi-Fi, background processes, sleep mode, and temperature can all affect the discharge curve.

For the best comparison between tests:

* start from 100%;
* use similar screen brightness;
* avoid heavy background tasks;
* keep the laptop awake during the test;
* stop the test at a similar battery percentage each time.

## Why this can be useful

This tool can help you compare:

* reported battery health vs measured discharge result;
* old battery vs replacement battery;
* different discharge tests on the same laptop;
* battery behavior after calibration.

I made it because I like old laptops and wanted a simple way to check battery capacity myself.

Maybe it will be useful for someone else too.

## License

MIT
