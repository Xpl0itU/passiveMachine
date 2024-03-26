package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
)

const (
	PAWNSAPP_IMAGE_NAME = "iproyal/pawns-cli:latest"
)

type PawnsAppConfig struct {
	Email      string
	Password   string
	DeviceName string
	Configured bool
}

func (i *PawnsAppConfig) ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application) {
	email := ""
	password := ""
	deviceName := ""
	isError := false
	showingError := false
	form.AddInputField("Device Name", i.DeviceName, 15, nil, func(text string) {
		deviceName = text
	})
	form.AddInputField("Email", i.Email, 50, nil, func(text string) {
		email = text
	})
	form.AddPasswordField("Password", i.Password, 20, '*', func(text string) {
		password = text
	})
	form.AddButton("Save", func() {
		isError = stringIsEmpty(email) || stringIsEmpty(password) || stringIsEmpty(deviceName)
		if isError {
			if !showingError {
				form.AddTextView("Error", "All fields are required", 0, 1, true, true)
				showingError = true
			}
			return
		}
		i.Email = email
		i.Password = password
		i.DeviceName = deviceName
		i.Configured = true
		returnToMenu(frame, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(frame, app)
	})
	form.AddButton("Register", func() {
		webbrowser.Open("https://pawns.app/?r=1112060")
	})
}

func (i *PawnsAppConfig) ConfigureDocker(kind DockerConfigKind, form *tview.Form) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `pawnsapp:
  image: ` + PAWNSAPP_IMAGE_NAME + `
  environment:
   - IPROYALPAWNS_DUMMY=''
  command: -accept-tos -email=` + i.Email + ` -password=` + i.Password + ` -device-name=` + i.DeviceName + ` -device-id=id_` + i.DeviceName + `
  restart: unless-stopped
`, nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		containerConfig := &container.Config{
			Image: PAWNSAPP_IMAGE_NAME,
			Cmd:   []string{"-accept-tos", "-email=" + i.Email, "-password=" + i.Password, "-device-name=" + i.DeviceName, "-device-id=id_" + i.DeviceName},
			Env: []string{
				"IPROYALPAWNS_DUMMY=",
			},
		}
		hostConfig := &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}
		return "", createContainer("pawnsapp", containerConfig, hostConfig, form)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *PawnsAppConfig) IsConfigured() bool {
	return i.Configured
}

func (i *PawnsAppConfig) PostConfigure(form *tview.Form, app *tview.Application) {
}
