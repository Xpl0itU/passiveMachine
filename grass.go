package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
)

const (
	GRASS_IMAGE_NAME = "mrcolorrain/grass:latest"
)

type GrassConfig struct {
	Email      string
	Password   string
	Configured bool
}

func (i GrassConfig) ConfigureForm(form *tview.Form, list *tview.List, app *tview.Application) {
	email := ""
	password := ""
	isError := false
	showingError := false
	form.AddInputField("Email", i.Email, 50, nil, func(text string) {
		email = text
	})
	form.AddPasswordField("Password", i.Password, 20, '*', func(text string) {
		password = text
	})
	form.AddButton("Save", func() {
		isError = stringIsEmpty(email) || stringIsEmpty(password)
		if isError {
			if !showingError {
				form.AddTextView("Error", "All fields are required", 0, 1, true, true)
				showingError = true
			}
			return
		}
		i.Email = email
		i.Password = password
		i.Configured = true
		returnToMenu(list, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(list, app)
	})
}

func (i GrassConfig) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `grass:
  image: ` + GRASS_IMAGE_NAME + `
  environment:
	- GRASS_USER=` + i.Email + `
	- GRASS_PASS=` + i.Password + `
  restart: unless-stopped
`, nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		containerConfig := &container.Config{
			Image: GRASS_IMAGE_NAME,
			Env: []string{
				"GRASS_USER=" + i.Email,
				"GRASS_PASS=" + i.Password,
			},
		}
		hostConfig := &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}
		return "", createContainer("grass", containerConfig, hostConfig, logView)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i GrassConfig) IsConfigured() bool {
	return i.Configured
}