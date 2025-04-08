package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/qeesung/image2ascii/convert"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Takes a pointer to a theater and gives it the default settings.
func setupTheater(t *Theater, r *lipgloss.Renderer, height int, width int) {
	t.fileNum = 1
	t.color = "#FFFFFF"
	t.opts = convert.DefaultOptions
	t.opts.Colored = false
	t.opts.FixedWidth = width
	t.opts.FixedHeight = height
	t.height = height
	t.width = width
	t.renderer = r
}

func main() {
	// setup theater struct
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	check(err)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH sever @ ", net.JoinHostPort(host, port))

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()
	<-done

	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

}

// used to handle running bubbletea over ssh
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s) // Removed unused variable

	t := new(Theater)
	setupTheater(t, renderer, pty.Window.Height, pty.Window.Width)
	return t, []tea.ProgramOption{tea.WithAltScreen()}
}
