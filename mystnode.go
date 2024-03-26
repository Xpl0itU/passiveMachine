package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
)

const (
	MYST_IMAGE_NAME    = "mysteriumnetwork/myst:latest"
	MYST_REFERRAL_LINK = "https://mystnodes.co/?referral_code=ijIy8nJv8xqVoshRmJjKATvoZZYKZ3jhzOY3FWy6"
)

type MystConfig struct {
	Configured bool
}

func (i *MystConfig) ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application) {
	enabled := i.Configured
	form.AddCheckbox("Enable Myst", i.Configured, func(checked bool) {
		enabled = checked
	})
	form.AddButton("Save", func() {
		i.Configured = enabled
		returnToMenu(frame, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(frame, app)
	})
	form.AddButton("Register", func() {
		modal := tview.NewModal().
			SetText("Register on Mysterium Nodes\n" + MYST_REFERRAL_LINK).
			AddButtons([]string{"Open", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Open" {
					webbrowser.Open(MYST_REFERRAL_LINK)
				}
				app.SetRoot(form, true)
			})
		app.SetRoot(modal, true)
	})
}

func (i *MystConfig) ConfigureDocker(kind DockerConfigKind, form *tview.Form) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `myst:
	image: ` + MYST_IMAGE_NAME + `
	environment:
		- MYSTNODE_DUMMY=''
	command: service --agreed-terms-and-conditions
	network_mode: host
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
			Volumes: map[string]struct{}{
				"/var/lib/mysterium-node": {},
			},
		}
		hostConfig := &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
			CapAdd:      []string{"NET_ADMIN"},
			NetworkMode: "host",
			PortBindings: map[nat.Port][]nat.PortBinding{
				nat.Port("4449/tcp"): {
					{
						HostPort: "4449",
					},
				},
			},
		}
		return "", createContainer("myst", containerConfig, hostConfig, form)
	}
	return "", errors.New("unknown kind")
}

func (i *MystConfig) IsConfigured() bool {
	return i.Configured
}

func (i *MystConfig) PostConfigure(form *tview.Form, app *tview.Application) {
	form.AddButton("Open Myst Node URL", func() {
		webbrowser.Open("http://127.0.0.1:4449/")
	})
}
