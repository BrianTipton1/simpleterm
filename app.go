package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"termvim/pkg/server"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go server.StartServer()
}

func (a *App) Print(s string) {
	fmt.Println(s)
}

func (a *App) RunCmd(s string) string {
	args := strings.Split(s, " ")

	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return string(stdout) + "\n" + err.Error()
	}

	return string(stdout)
}
