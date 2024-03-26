package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
)

const (
	PACKETSTREAM_IMAGE_NAME = "packetstream/psclient:latest"
)

type PacketStreamConfig struct {
	CID        string
	Configured bool
}

func (i *PacketStreamConfig) ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application) {
	cid := ""
	isError := false
	showingError := false
	form.AddInputField("CID", i.CID, 50, nil, func(text string) {
		cid = text
	})
	form.AddButton("Save", func() {
		isError = stringIsEmpty(cid)
		if isError {
			if !showingError {
				form.AddTextView("Error", "All fields are required", 0, 1, true, true)
				showingError = true
			}
			return
		}
		i.CID = cid
		i.Configured = true
		returnToMenu(frame, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(frame, app)
	})
	form.AddButton("Register", func() {
		webbrowser.Open("https://packetstream.io/?psr=4cRE")
	})
}

func (i *PacketStreamConfig) ConfigureDocker(kind DockerConfigKind, form *tview.Form) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `packetstream:
  image: ` + PACKETSTREAM_IMAGE_NAME + `
  environment:
	- CID=` + i.CID + `
  restart: unless-stopped
`, nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		containerConfig := &container.Config{
			Image: PACKETSTREAM_IMAGE_NAME,
			Env: []string{
				"CID=" + i.CID,
			},
		}
		hostConfig := &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}
		return "", createContainer("packetstream", containerConfig, hostConfig, form)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *PacketStreamConfig) IsConfigured() bool {
	return i.Configured
}

func (i *PacketStreamConfig) PostConfigure(form *tview.Form, app *tview.Application) {
}
