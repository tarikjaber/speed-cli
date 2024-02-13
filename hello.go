package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	taskTime = 1000
)

type model struct {
	progress progress.Model
	pastScores []int
	score int
	tasksCompleted int
}

type tickMsg time.Time

func initialModel() model {
	progress := progress.New(
		progress.WithoutPercentage(),
		progress.WithWidth(50),
	)
	return model {
		progress: progress,
		pastScores: []int{},
		score: taskTime,
		tasksCompleted: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), tea.ClearScreen)
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
			m.score = taskTime
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
	s += fmt.Sprintf("Completed: %d\n\n", m.tasksCompleted)
	s += m.progress.ViewAs(math.Max(0, float64(averageScore) / taskTime)) + "\n"
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