package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qeesung/image2ascii/convert"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Takes a pointer to a theater and gives it the default settings.
func setupTheater(t *Theater) {
	t.fileNum = 1
	t.opts = convert.DefaultOptions
	t.opts.StretchedScreen = true // stretch to fit
}

func main() {
	// setup theater struct
	t := new(Theater)
	setupTheater(t)

	p := tea.NewProgram(t)

	_, err := p.Run()
	check(err)

}
