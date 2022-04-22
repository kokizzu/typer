package typer

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/typer/pkg/model"
	util "github.com/maaslalani/typer/pkg/utility"
	wrap "github.com/mitchellh/go-wordwrap"
)

const (
	blue          = "#4776E6"
	words         = 15
	defaultWidth  = 60
	DefaultLength = 20
	MaxLength     = 500
)

type Flags struct {
	Length      int
	Capital     bool
	Punctuation bool
}

// FromStdin takes input text from stdin
func FromStdin(n int, flagStruct *Flags) error {
	var stdin []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		stdin = append(stdin, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	text := string(stdin)
	text, err := flagStruct.formatText(text)
	if err != nil {
		return err
	}

	return run(text)
}

// FromRandom runs the app with text generated by RandomWords
func FromRandom(n int, flagStruct *Flags) error {
	text := util.RandomWords(n)

	text, err := flagStruct.formatText(text)
	if err != nil {
		return err
	}

	return run(text)
}

// FromFile runs the app with text extracted from a file
func FromFile(path string, flagStruct *Flags) error {
	text, err := util.ReadFile(path)
	if err != nil {
		return err
	}

	text, err = flagStruct.formatText(text)
	if err != nil {
		return err
	}

	return run(text)
}

// run is responsible for running the GUI
func run(text string) error {
	bar, err := progress.NewModel(progress.WithSolidFill(blue))
	if err != nil {
		return err
	}

	program := tea.NewProgram(model.Model{
		Progress: bar,
		Text:     wrap.WrapString(text, defaultWidth),
	})

	return program.Start()
}

// formatText applies formatting based on flags
func (f *Flags) formatText(s string) (string, error) {
	var err error
	s, err = util.AdjustWhitespace(s)
	if err != nil {
		return "", err
	}

	if !f.Punctuation {
		s, err = util.RemoveNonAlpha(s)
		if err != nil {
			return "", err
		}
	}

	if f.Length <= 0 {
		log.Println("Length value is incorrect. Restoring to default value.")
		f.Length = DefaultLength
	} else if f.Length > MaxLength {
		log.Println("Max length value exceeded. Restoring to max length value.")
		f.Length = MaxLength
	}
	s = util.AdjustLength(s, f.Length)

	if !f.Capital {
		s = strings.ToLower(s)
	}
	return s, nil
}
