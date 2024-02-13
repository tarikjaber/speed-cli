package main

import (
	"fmt"
	"os"
	"time"

	// "github.com/charmbracelet/bubbles/progress"
	// "github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	pastScores []int
	score int
	tasksCompleted int
}

type tickMsg time.Time

func initialModel() model {
	return model {
		pastScores: []int{},
		score: 1000,
		tasksCompleted: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update (msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n":
			m.tasksCompleted += 1
			m.pastScores = append(m.pastScores, m.score)
			m.score = 1000
		}
	case tickMsg:
		m.score -= 1
		return m, tickCmd()
	}
	var cmd tea.Cmd
	return m, cmd
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second * 1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) View() string {
	sumScores := m.score
	for _, score := range m.pastScores {
		sumScores += score
	}
	
	averageScore := sumScores / (m.tasksCompleted + 1)

	s := fmt.Sprintf("Average: %d\n", averageScore)
	s += fmt.Sprintf("Score: %d\n", m.score)
	s += fmt.Sprintf("Completed: %d\n", m.tasksCompleted)
	s += "\nPress n to continue.\n"
	s += "Press q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There was an error: %v", err)
		os.Exit(1)
	}
}