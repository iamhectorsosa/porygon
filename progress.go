package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/iamhectorsosa/porygon/internal/mapper"
)

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	dim    = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

type mapDoneMsg struct {
	res *mapper.Result
	err error
}

type model struct {
	spinner spinner.Model
	opts    mapper.Options
	res     *mapper.Result
	err     error
	done    bool
}

func newModel(opts mapper.Options) model {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, opts: opts}
}

func mapCmd(opts mapper.Options) tea.Cmd {
	return func() tea.Msg {
		res, err := mapper.Map(opts)
		return mapDoneMsg{res: res, err: err}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, mapCmd(m.opts))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case mapDoneMsg:
		m.res, m.err, m.done = msg.res, msg.err, true
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.err, m.done = errors.New("interrupted"), true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	head := fmt.Sprintf("🔺 %s is sniffing around %s",
		accent.Render("Porygon"), dim.Render(displayDir(m.opts.Dir)))

	if !m.done {
		return head + " " + m.spinner.View() + "\n"
	}
	if m.err != nil {
		return head + "\n"
	}
	return head + "\n" + summary(m.res) + "\n"
}

func displayDir(dir string) string {
	if dir == "." {
		return "the current directory"
	}
	return dir
}

func summary(res *mapper.Result) string {
	dirs := 0
	for _, e := range res.Entries {
		if e.IsDir {
			dirs++
		}
	}
	files := len(res.Entries) - dirs

	return fmt.Sprintf("🔺 %s mapped the directory. It's super effective! %s",
		accent.Render("Porygon"),
		dim.Render(fmt.Sprintf("%d folders · %d files · %d read", dirs, files, len(res.Files))),
	)
}
