package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/qeesung/image2ascii/convert"
)

// struct that represents the frame
type frame string

// bubbletea model that will play badapple
type Theater struct {
	fileNum  int
	opts     convert.Options
	frameStr string
}

func (t Theater) Init() tea.Cmd {
	f, err := tea.LogToFile("debug.log", "debug")
	log.SetLevel(log.DebugLevel)
	log.SetOutput(f)
	if err != nil {

	}
	defer f.Close()
	return t.showFrame
}

func (t Theater) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			log.Debug("Exit on user request")
			fmt.Println(t.frameStr)
			return t, tea.Quit
		}
	case frame:
		log.Debug("Recieved Frame")
		if t.fileNum <= 6573 {
			newT := t
			newT.fileNum++
			newT.frameStr = string(msg)
			return newT, newT.showFrame
		}
	}

	return t, nil
}

func (t Theater) View() string {
	return t.frameStr
}

func (t Theater) showFrame() tea.Msg {
	filePath := fmt.Sprintf("./frames/%04d.jpg", t.fileNum)

	c := convert.NewImageConverter()

	t.frameStr = c.ImageFile2ASCIIString(filePath, &t.opts)

	return frame(t.frameStr)
}
