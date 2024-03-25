package main

import "github.com/rivo/tview"

type MenuItemConfig interface {
	ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application)
	ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error)
	IsConfigured() bool
}

type MenuItem struct {
	Name        string
	Description string
	Config      MenuItemConfig
}

func (i *MenuItem) GetName() string {
	return i.Name
}

func (i *MenuItem) GetDescription() string {
	return i.Description
}

func returnToMenu(frame *tview.Frame, app *tview.Application) {
	app.SetRoot(frame, true)
}
