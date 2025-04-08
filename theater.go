package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/qeesung/image2ascii/convert"
)

// struct that represents the frame
type frame string

type tickMsg time.Time

// bubbletea model that will play badapple
type Theater struct {
	fileNum  int
	opts     convert.Options
	frameStr string
	color    string
	renderer *lipgloss.Renderer
	width    int
	height   int
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
			return t, tea.Quit
		// SECTION: COLOR CHANGES
		case "`":
			// White
			t.color = "#FFFFFF"
		case "1":
			// Red
			t.color = "#FF191C"
		case "2":
			// blue
			t.color = "#0F33FF"
		case "3":
			// green
			t.color = "#4CFF4D"
		}
	case tickMsg:
		if t.fileNum <= 6573 {
			newT := t
			newT.fileNum++
			return newT, newT.showFrame
		}

	case frame:
		log.Debug("Recieved Frame")
		if t.fileNum <= 6573 {
			newT := t
			newT.fileNum++
			newT.frameStr = string(msg)
			return newT, tick()
		}
	case tea.WindowSizeMsg:
		t.height = msg.Height
		t.width = msg.Width
	}

	return t, nil
}

func (t Theater) View() string {
	style := t.renderer.NewStyle().SetString(t.frameStr).Foreground((lipgloss.Color(t.color)))

	return style.Render()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (t Theater) showFrame() tea.Msg {
	frameDuration := time.Second / 30
	time.Sleep(frameDuration)
	filePath := fmt.Sprintf("./frames/%04d.jpg", t.fileNum)

	c := convert.NewImageConverter()

	t.frameStr = c.ImageFile2ASCIIString(filePath, &t.opts)

	return frame(t.frameStr)
}
