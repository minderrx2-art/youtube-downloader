package internal

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type VideoDebug struct {
	id       int
	progress string
	message  string
}

type VideoState struct {
	id       int
	progress float64 // 0.0 -> 1.0
	title    string
	log      []string
}

type model struct {
	progressBar progress.Model
	videos      map[int]VideoState
}

func initStatus(id int, title string) tea.Cmd {
	return func() tea.Msg {
		return VideoState{
			id:    id,
			title: title,
			log:   make([]string, 0),
		}
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for id, state := range m.videos {
		cmds = append(cmds, initStatus(id, state.title))
	}

	return tea.Batch(cmds...)
}

func (m model) allComplete() bool {
	for _, state := range m.videos {
		if state.progress < 1.0 {
			return false
		}
	}

	return true
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case VideoDebug:
		state := m.videos[msg.id]

		if msg.message != "" {
			state := m.videos[msg.id]
			state.log = append(state.log, msg.message)
		}

		pct, _ := strconv.ParseFloat(msg.progress, 64)
		state.progress = pct / 100.0

		m.videos[msg.id] = state

		if m.allComplete() {
			return m, tea.Quit
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	var out strings.Builder

	for i := 0; i < len(m.videos); i++ {
		v := m.videos[i]

		out.WriteString(v.title)
		out.WriteString("\n")

		out.WriteString(
			m.progressBar.ViewAs(v.progress),
		)
		out.WriteString("\n\n")
	}

	return out.String()
}

func NewOutput(titles []string, size int) model {
	videos := make(map[int]VideoState, size)

	bar := progress.New()
	for i := 0; i < size; i++ {
		videos[i] = VideoState{
			id:       i,
			progress: 0,
			title:    titles[i],
		}
	}

	return model{
		progressBar: bar,
		videos:      videos,
	}
}
