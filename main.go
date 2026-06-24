package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/CtrlPy/batcap/battery"
	"github.com/CtrlPy/batcap/report"
	"github.com/CtrlPy/batcap/tui"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"
)

func main() {
	reset := flag.BoolP("reset", "r", false, "Reset previous session")
	batt := flag.StringP("battery", "b", "BAT0", "Battery name (e.g., BAT0)")
	flag.Parse()

	reader := battery.NewReader(*batt)
	info, err := reader.ReadInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading battery info: %v\n", err)
		os.Exit(1)
	}

	if info.Status == "Charging" || info.Status == "Full" {
		fmt.Println("WARNING: The battery is currently charging or full. Please disconnect the AC adapter to discharge.")
	}

	if info.Capacity < 95 && !*reset {
		fmt.Println("WARNING: For accurate results, it is highly recommended to start from 100% capacity.")
	}

	sess, err := battery.NewSession(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize session: %v\n", err)
		os.Exit(1)
	}

	if err := sess.StartOrResume(*reset); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start session: %v\n", err)
		os.Exit(1)
	}

	// Handle graceful shutdown on OS signals (e.g., if system starts shutting down)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-sigs
		sess.Stop()
		saveReport(sess)
		os.Exit(0)
	}()

	p := tea.NewProgram(tui.NewModel(sess))

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	sess.Stop()
	saveReport(sess)
}

func saveReport(sess *battery.Session) {
	// Print final report to console
	finalReport := report.Generate(sess.State)
	fmt.Println()
	fmt.Println(finalReport)

	// Save to file in the current directory
	cwd, err := os.Getwd()
	if err == nil {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := filepath.Join(cwd, fmt.Sprintf("batcap_report_%s.txt", timestamp))

		// Remove formatting codes for the plain text file
		fileContent := "BATCAP REPORT - " + timestamp + "\n\n" + stripANSI(finalReport) + "\n"

		if err := os.WriteFile(filename, []byte(fileContent), 0644); err == nil {
			fmt.Printf("\n[+] Report successfully saved to: %s\n", filename)
		} else {
			fmt.Printf("\n[-] Failed to save report to file: %v\n", err)
		}
	}
}

// Simple ANSI stripper for the plain text file
func stripANSI(str string) string {
	var sb []rune
	inEscape := false
	for _, r := range str {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		sb = append(sb, r)
	}
	return string(sb)
}
