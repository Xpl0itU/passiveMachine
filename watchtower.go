package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
)

const (
	WATCHTOWER_IMAGE_NAME = "containrrr/watchtower:latest"
)

type WatchtowerConfig struct {
	Configured bool
}

func (i *WatchtowerConfig) ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application) {
	enabled := i.Configured
	form.AddCheckbox("Automatic Updates", enabled, func(checked bool) {
		enabled = checked
	})
	form.AddButton("Save", func() {
		i.Configured = enabled
		returnToMenu(frame, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(frame, app)
	})
}

func (i *WatchtowerConfig) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `watchtower:
  image: ` + WATCHTOWER_IMAGE_NAME + `
  volumes:
	- /var/run/docker.sock:/var/run/docker.sock
  restart: always
`, nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		containerConfig := &container.Config{
			Image: WATCHTOWER_IMAGE_NAME,
		}
		hostConfig := &container.HostConfig{
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		}
		return "", createContainer("watchtower", containerConfig, hostConfig, logView)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *WatchtowerConfig) IsConfigured() bool {
	return i.Configured
}
