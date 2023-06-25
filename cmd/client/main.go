package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	mode      mode
	conn      net.Conn
	textInput textinput.Model

	highlight  int
	characters characters
}

type character struct {
	name       string
	initiative string
}
type characters []character

func (c characters) Len() int           { return len(c) }
func (c characters) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c characters) Less(i, j int) bool { return c[i].initiative > c[j].initiative }

//go:generate stringer -type=mode
type mode uint8

const (
	initiative mode = iota
	battle
	input
)

var (
	textStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("7")).
			PaddingRight(2)

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6")).
			PaddingRight(2)

	headingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("7")).
			Background(lipgloss.Color("6"))

	blockStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingBottom(1)
)

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case initiative:
		return m.initiativeUpdate(msg)
	case battle:
		return m.battleUpdate(msg)
	case input:
		return m.inputUpdate(msg)
	default:
		return m, nil
	}
}

func (m model) initiativeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j":
			m.highlight += 1
			if m.highlight > len(m.characters)-1 {
				m.highlight = 0
			}
		case "k":
			m.highlight -= 1
			if m.highlight < 0 {
				m.highlight = len(m.characters) - 1
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			switch len(m.characters[m.highlight].initiative) {
			case 0:
				m.characters[m.highlight].initiative += msg.String()
			case 1:
				m.characters[m.highlight].initiative += msg.String()
				m.highlight += 1
				if m.highlight > len(m.characters)-1 {
					m.highlight = 0
				}
			default:
				m.characters[m.highlight].initiative = msg.String()
			}
		case "x":
			if len(m.characters) == 0 {
				break
			}
			m.characters = append(
				m.characters[:m.highlight],
				m.characters[m.highlight+1:]...,
			)
			if m.highlight > len(m.characters)-1 {
				m.highlight = 0
			}
		case "n":
			m.mode = input
		case "enter":
			sort.Sort(m.characters)
			var b strings.Builder
			b.WriteString("battle")
			for _, c := range m.characters {
				b.WriteString(",")
				b.WriteString(c.name)
			}
			fmt.Fprintln(m.conn, b.String())
			m.mode = battle
		}
	}
	return m, nil
}

func (m model) battleUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			m.mode = initiative
		case "j", "enter", "space":
			m.highlight += 1
			if m.highlight > len(m.characters)-1 {
				m.highlight = 0
			}
			fmt.Fprintf(m.conn, "highlight,%d\n", m.highlight)
		case "k", "backspace":
			m.highlight -= 1
			if m.highlight < 0 {
				m.highlight = len(m.characters) - 1
			}
			fmt.Fprintf(m.conn, "highlight,%d\n", m.highlight)
		}
	}
	return m, nil
}

func (m model) inputUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			m.mode = initiative
		case "enter":
			c := character{
				name: m.textInput.Value(),
			}
			m.characters = append(m.characters, c)
			m.mode = initiative
		default:
			m.textInput, cmd = m.textInput.Update(msg)
		}
	}
	return m, cmd
}

func (m model) View() string {
	var names []string
	for i, c := range m.characters {
		if i == m.highlight {
			names = append(names, highlightStyle.Render(c.name))
		} else {
			names = append(names, textStyle.Render(c.name))
		}
	}
	namesRendered := lipgloss.JoinVertical(lipgloss.Left, names...)

	var initiatives []string
	for i, c := range m.characters {
		if i == m.highlight {
			initiatives = append(initiatives, highlightStyle.Render(c.initiative))
		} else {
			initiatives = append(initiatives, textStyle.Render(c.initiative))
		}
	}
	initiativesRendered := lipgloss.JoinVertical(lipgloss.Left, initiatives...)

	title := headingStyle.Render(m.mode.String())
	block := blockStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		namesRendered,
		initiativesRendered,
	))

	if m.mode != input {
		return lipgloss.JoinVertical(lipgloss.Left, title, block)
	}
	return lipgloss.JoinVertical(lipgloss.Left, title, block, m.textInput.View())
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading config: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", config.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed connecting to server: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(conn, "start")

	var characters []character
	for _, name := range config.Party {
		characters = append(characters, character{name: name})
	}

	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := model{
		characters: characters,
		conn:       conn,
		textInput:  ti,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting initiative client: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(conn, "end")
}
