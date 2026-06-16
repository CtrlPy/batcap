package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/CtrlPy/batcap/battery"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)
			
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E1ED"))
)

type tickMsg time.Time
type autoStopMsg struct{}

type Model struct {
	sess *battery.Session
	err  error
}

func NewModel(sess *battery.Session) Model {
	return Model{sess: sess}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		m.listenForAutoStop(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) listenForAutoStop() tea.Cmd {
	return func() tea.Msg {
		<-m.sess.AutoStop
		return autoStopMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		return m, tickCmd()
	case autoStopMsg:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	state := m.sess.State
	power := state.LastPower

	var sb strings.Builder

	sb.WriteString(titleStyle.Render(" batcap — Live Monitoring ") + "\n")

	bmsHealth := 0.0
	if state.EnergyFullDesign > 0 {
		bmsHealth = (state.EnergyFull / state.EnergyFullDesign) * 100.0
	}
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(fmt.Sprintf("Battery: %s | BMS Health: %.0f%% | Cycles: %d", state.BatteryModel, bmsHealth, state.CycleCount)) + "\n\n")

	sb.WriteString(labelStyle.Render("Elapsed Time:       ") + valueStyle.Render(battery.FormatDuration(state.ElapsedSeconds)) + "\n")
	sb.WriteString(labelStyle.Render("Current Power Draw: ") + valueStyle.Render(fmt.Sprintf("%.2f W", power)) + "\n")
	sb.WriteString(labelStyle.Render("Battery Capacity:   ") + valueStyle.Render(fmt.Sprintf("%d%%", state.CurrentCapacity)) + "\n")
	
	sb.WriteString("\n")
	sb.WriteString(labelStyle.Render("Real Energy Discharged (Integral): ") + valueStyle.Render(fmt.Sprintf("%.2f Wh", state.IntegratedEnergy)) + "\n")
	bmsDiff := state.EnergyStart - state.EnergyCurrentBMS
	sb.WriteString(labelStyle.Render("BMS Estimated Discharged (Diff):   ") + valueStyle.Render(fmt.Sprintf("%.2f Wh", bmsDiff)) + "\n")
	
	sb.WriteString("\n")
	
	if len(state.PowerHistory) > 1 {
		sb.WriteString(labelStyle.Render("Power Draw History (last 60s):") + "\n")
		graph := asciigraph.Plot(state.PowerHistory, asciigraph.Height(5), asciigraph.Width(50))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E1ED")).Render(graph) + "\n\n")
	}

	sb.WriteString(lipgloss.NewStyle().Faint(true).Render("Press 'q' or Ctrl+C to stop and generate report."))

	return lipgloss.NewStyle().Margin(1, 2).Render(sb.String())
}
