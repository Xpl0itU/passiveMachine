package main

import "github.com/rivo/tview"

type MenuItem interface {
	GetName() string
	GetDescription() string
	ConfigureForm(form *tview.Form, list *tview.List, app *tview.Application)
	ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error)
	IsConfigured() bool
}

func returnToMenu(list *tview.List, app *tview.Application) {
	app.SetRoot(list, true)
}
