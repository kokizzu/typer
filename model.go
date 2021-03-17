package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	width = 60.

	// charsPerWord is the average characters per word used by most typing tests
	// to calculate your WPM score.
	charsPerWord = 5.
)

type model struct {
	// percent is a value from 0 to 1 that represents the current completion of the typing test
	percent  float64
	progress *progress.Model
	// text is the randomly generated text for the user to type
	text string
	// typed is the text that the user has typed so far
	typed string
	// start and end are the start and end time of the typing test
	start time.Time
	end   time.Time
	// mistakes is the number of characters that were mistyped by the user
	mistakes int
	// score is the user's score calculated by correct characters typed
	score float64
}

// Init inits the bubbletea model for use
func (m model) Init() tea.Cmd {
	return nil
}

// Update updates the bubbletea model by handling the progress bar update
// and adding typed characters to the state if they are valid typing characters
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// User wants to cancel the typing test
		if msg.Type == tea.KeyCtrlC {
			m.end = time.Now()
			return m, tea.Quit
		}

		// Deleting characters
		if msg.Type == tea.KeyBackspace && len(m.typed) > 0 {
			m.typed = m.typed[:len(m.typed)-1]
		}

		// Ensure we are adding characters only that we want the user to be able to type
		// and that all the typed characters have a rune width of one.
		if len(msg.Runes) > 0 && int(msg.Runes[0]) > 31 && int(msg.Runes[0]) < 127 {
			m.typed += msg.String()
		}

		// Update progress bar state
		m.percent = float64(len(m.typed)) / float64(len(m.text))
		if m.percent >= 1.0 {
			m.end = time.Now()
			return m, tea.Quit
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		if m.progress.Width > width {
			m.progress.Width = width
		}
		return m, nil

	default:
		return m, nil
	}
}

// View shows the current state of the typing test.
// It displays a progress bar for the progression of the typing test,
// the typed characters (with errors displayed in red) and remaining
// characters to be typed in a faint display
func (m model) View() string {
	remaining := m.text[len(m.typed):]

	var typed string
	for i, c := range m.typed {
		if byte(c) == m.text[i] {
			typed += string(c)
			m.score += 1.0
			continue
		}
		typed += red(string(c))
	}

	s := fmt.Sprintf("\n  %s\n\n%s%s", m.progress.View(m.percent), typed, faint(remaining))

	// Display words per minute when finished
	if len(m.typed) >= len(m.text) {
		s += fmt.Sprintf(bold("\n\nWPM: %.2f\n"), (m.score/charsPerWord)/(m.end.Sub(m.start).Minutes()))
	}
	return s
}