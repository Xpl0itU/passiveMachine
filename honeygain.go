package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
)

const (
	HONEYGAIN_IMAGE_NAME    = "honeygain/honeygain:latest"
	HONEYGAIN_REFERRAL_LINK = "https://r.honeygain.me/SAMUEC73"
)

type HoneygainConfig struct {
	DeviceName string
	Email      string
	Password   string
	Configured bool
}

func (i *HoneygainConfig) ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application) {
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
		modal := tview.NewModal().
			SetText("Register on Honeygain\n" + HONEYGAIN_REFERRAL_LINK).
			AddButtons([]string{"Open", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Open" {
					webbrowser.Open(HONEYGAIN_REFERRAL_LINK)
				}
				app.SetRoot(form, true)
			})
		app.SetRoot(modal, true)
	})
}

func (i *HoneygainConfig) ConfigureDocker(kind DockerConfigKind, form *tview.Form) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `honeygain:
	image: ` + HONEYGAIN_IMAGE_NAME + `
	restart: unless-stopped
	environment:
		- HONEYGAIN_DUMMY=''
	command: -tou-accept -email ` + i.Email + ` -pass ` + i.Password + ` -device ` + i.DeviceName + "\n", nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		hostConfig := &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}
		containerConfig := &container.Config{
			Image: HONEYGAIN_IMAGE_NAME,
			Env: []string{
				"HONEYGAIN_DUMMY=",
			},
			Cmd: []string{
				"-tou-accept",
				"-email",
				i.Email,
				"-pass",
				i.Password,
				"-device",
				i.DeviceName,
			},
		}
		return "", createContainer("honeygain", containerConfig, hostConfig, form)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *HoneygainConfig) IsConfigured() bool {
	return i.Configured
}

func (i *HoneygainConfig) PostConfigure(form *tview.Form, app *tview.Application) {
}
