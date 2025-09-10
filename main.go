package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

const (
	host = "localhost"
	port = "6942"
)

type app struct {
	server *ssh.Server
	progs  []*tea.Program
}

func (a app) send(msg tea.Msg) {
	for _, p := range a.progs {
		go p.Send(msg)
	}
}

func newApp() *app {
	a := new(app)
	s, err := wish.NewServer(wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(bubbletea.MiddlewareWithProgramHandler(a.ProgramHandler, termenv.ANSI256),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}
	a.server = s
	return a
}

func (a *app) Start() {
	var err error
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = a.server.ListenAndServe(); err != nil {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}

func (a *app) ProgramHandler(s ssh.Session) *tea.Program {
	var dump *os.File
	os.Setenv("DEBUG", "true")
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	model, _ := initialBaseModel(dump)

	p := tea.NewProgram(model, bubbletea.MakeOptions(s)...)
	a.progs = append(a.progs, p)

	return p
}

func main() {
	app := newApp()
	app.Start()
}

//
// func main() {
// 	// var dump *os.File
// 	// a := 1
// 	// os.Setenv("DEBUG", "true")
// 	// if _, ok := os.LookupEnv("DEBUG"); ok {
// 	// 	a = 2
// 	// 	var err error
// 	// 	dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
// 	// 	if err != nil {
// 	// 		os.Exit(1)
// 	// 	}
// 	// }
// 	// println(a)
// 	// m, err := initialBaseModel(dump)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// p := tea.NewProgram(m, tea.WithAltScreen())
// 	// // p := tea.NewProgram(initialList(), tea.WithAltScreen())
// 	// go func() {
// 	// 	checkFolderUpdates(p)
// 	// }()
// 	// _, err = p.Run()
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	s, err := wish.NewServer(
// 		wish.WithAddress(net.JoinHostPort(host, port)),
// 		wish.WithHostKeyPath(".ssh/id_ed25519"),
// 		wish.WithMiddleware(
// 			bubbletea.Middleware(teaHandler),
// 			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
// 			logging.Middleware(),
// 		),
// 	)
// 	if err != nil {
// 		log.Error("Could not start server", "error", err)
// 	}
//
// 	done := make(chan os.Signal, 1)
// 	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
// 	log.Info("Starting SSH server", "host", host, "port", port)
// 	go func() {
// 		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
// 			log.Error("Could not start server", "error", err)
// 			done <- nil
// 		}
// 	}()
//
// 	<-done
// 	log.Info("Stopping SSH server")
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer func() { cancel() }()
// 	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
// 		log.Error("Could not stop server", "error", err)
// 	}
// }
//
// func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
// 	// This should never fail, as we are using the activeterm middleware.
// 	s.Pty()
//
// 	// When running a Bubble Tea app over SSH, you shouldn't use the default
// 	// lipgloss.NewStyle function.
// 	// That function will use the color profile from the os.Stdin, which is the
// 	// server, not the client.
// 	// We provide a MakeRenderer function in the bubbletea middleware package,
// 	// so you can easily get the correct renderer for the current session, and
// 	// use it to create the styles.
// 	// The recommended way to use these styles is to then pass them down to
// 	// your Bubble Tea model.
//
// 	var dump *os.File
// 	a := 1
// 	os.Setenv("DEBUG", "true")
// 	if _, ok := os.LookupEnv("DEBUG"); ok {
// 		a = 2
// 		var err error
// 		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
// 		if err != nil {
// 			os.Exit(1)
// 		}
// 	}
// 	println(a)
// 	m, err := initialBaseModel(dump)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return m, []tea.ProgramOption{tea.WithAltScreen()}
// }
