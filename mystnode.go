package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/rivo/tview"
)

const (
	MYST_IMAGE_NAME = "mysteriumnetwork/myst:latest"
)

type MystConfig struct {
	Configured bool
}

func (i *MystConfig) ConfigureForm(form *tview.Form, list *tview.List, app *tview.Application) {
	enabled := i.Configured
	form.AddCheckbox("Enable Myst", i.Configured, func(checked bool) {
		enabled = checked
	})
	form.AddButton("Save", func() {
		i.Configured = enabled
		returnToMenu(list, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(list, app)
	})
}

func (i *MystConfig) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `myst:
  image: ` + MYST_IMAGE_NAME + `
  environment:
    - MYSTNODE_DUMMY=''
  command: service --agreed-terms-and-conditions
  cap_add:
    - NET_ADMIN
  ports:
    - "4449:4449"
  volumes:
    - myst-data:/var/lib/mysterium-node
  restart: unless-stopped
`, nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		containerConfig := &container.Config{
			Image: MYST_IMAGE_NAME,
			Env: []string{
				"MYSTNODE_DUMMY=",
			},
			Cmd: []string{"service", "--agreed-terms-and-conditions"},
		}
		hostConfig := &container.HostConfig{
			VolumesFrom: []string{"myst-data"},
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
			CapAdd: []string{"NET_ADMIN"},
			PortBindings: map[nat.Port][]nat.PortBinding{
				"4449/tcp": {
					{
						HostIP:   "0.0.0.0",
						HostPort: "4449",
					},
				},
			},
		}
		return "", createContainer("myst", containerConfig, hostConfig, logView)
	}
	return "", errors.New("unknown kind")
}

func (i *MystConfig) IsConfigured() bool {
	return i.Configured
}
