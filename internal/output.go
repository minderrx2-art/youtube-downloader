package internal

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type VideoDebug struct {
	id      int
	done    bool
	message string
}

type VideoState struct {
	id   int
	done bool
	view string
}

type model struct {
	progress map[int]VideoState
}

func initStatus(id int, title string) tea.Cmd {
	return func() tea.Msg {
		return VideoState{
			id:   id,
			view: title,
		}
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for id, state := range m.progress {
		cmds = append(cmds, initStatus(id, state.view))
	}

	return tea.Batch(cmds...)
}

func (m model) allComplete() bool {
	for _, state := range m.progress {
		if state.done == false {
			return false
		}
	}

	return true
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case VideoDebug:
		state := m.progress[msg.id]
		state.view = msg.message
		m.progress[msg.id] = state

		if msg.done == true {
			state.done = true
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

	for i := 0; i < len(m.progress); i++ {
		out.WriteString(m.progress[i].view)
		out.WriteString("\n")
	}

	return out.String()
}

func NewOutput(titles []string, size int) model {
	progress := make(map[int]VideoState, size)

	for i := 0; i < size; i++ {
		progress[i] = VideoState{
			id:   i,
			done: false,
			view: titles[i],
		}
	}

	return model{
		progress: progress,
	}
}
