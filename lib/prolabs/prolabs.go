package prolabs

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel() model {
	return model{
		Prolabs: []Prolabs{
			{"Dante", 10, 5, "Intermediate", 20},
			{"Offshore", 15, 8, "Advanced", 50},
			{"Zephyr", 15, 8, "Intermediate", 50},
			{"Rastalabs", 15, 8, "Intermediate", 50},
			{"Cybernetics", 15, 8, "Advanced", 50},
			{"APTLabs", 15, 8, "Advanced", 50},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var leftColumn, rightColumn strings.Builder
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(m.width/2 - 4)
	currentHeight := 0
	maxHeight := m.height

	for _, prolabs := range m.Prolabs {
		block := boxStyle.Render(
			fmt.Sprintf("Nom: %s\nFlags: %d\nMachines: %d\nDifficulté: %s\nProgression: %d%%",
				prolabs.Name,
				prolabs.Flags,
				prolabs.Machines,
				prolabs.Difficulty,
				prolabs.Progression,
			),
		)

		blockHeight := len(strings.Split(block, "\n")) + 1

		if currentHeight+blockHeight > maxHeight {
			rightColumn.WriteString(block)
			rightColumn.WriteString("\n")
			currentHeight = blockHeight
		} else {
			leftColumn.WriteString(block)
			leftColumn.WriteString("\n")
			currentHeight += blockHeight
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn.String(), rightColumn.String())
}
