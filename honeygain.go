package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
)

const (
	HONEYGAIN_IMAGE_NAME = "honeygain/honeygain:latest"
)

type HoneygainConfig struct {
	DeviceName string
	Email      string
	Password   string
}

type HoneygainItem struct {
	Name        string
	Description string
	Config      HoneygainConfig
	Configured  bool
}

func (i *HoneygainItem) GetName() string {
	return i.Name
}

func (i *HoneygainItem) GetDescription() string {
	return i.Description
}

func (i *HoneygainItem) ConfigureForm(form *tview.Form, list *tview.List, app *tview.Application) {
	email := ""
	password := ""
	deviceName := ""
	isError := false
	form.AddInputField("Device Name", i.Config.DeviceName, 15, nil, func(text string) {
		deviceName = text
	})
	form.AddInputField("Email", i.Config.Email, 50, nil, func(text string) {
		email = text
	})
	form.AddPasswordField("Password", i.Config.Password, 20, '*', func(text string) {
		password = text
	})
	form.AddButton("Save", func() {
		if !isError && (stringIsEmpty(email) || stringIsEmpty(password) || stringIsEmpty(deviceName)) {
			form.AddTextView("Error", "All fields are required", 0, 1, true, true)
			isError = true
			return
		}
		i.Config.Email = email
		i.Config.Password = password
		i.Config.DeviceName = deviceName
		i.Configured = true
		returnToMenu(list, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(list, app)
	})
}

func (i *HoneygainItem) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `honeygain:
  image: ` + HONEYGAIN_IMAGE_NAME + `
  restart: unless-stopped
  environment:
	- HONEYGAIN_DUMMY=''
  command: -tou-accept -email ` + i.Config.Email + ` -pass ` + i.Config.Password + ` -device ` + i.Config.DeviceName + "\n", nil
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
				i.Config.Email,
				"-pass",
				i.Config.Password,
				"-device",
				i.Config.DeviceName,
			},
		}
		return "", createContainer("honeygain", containerConfig, hostConfig, logView)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *HoneygainItem) IsConfigured() bool {
	return i.Configured
}
