package main

import (
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
)

const (
	PACKETSTREAM_IMAGE_NAME = "packetstream/psclient:latest"
)

type PacketStreamConfig struct {
	CID        string
	Configured bool
}

func (i *PacketStreamConfig) ConfigureForm(form *tview.Form, list *tview.List, app *tview.Application) {
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
		returnToMenu(list, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(list, app)
	})
}

func (i *PacketStreamConfig) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
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
		return "", createContainer("packetstream", containerConfig, hostConfig, logView)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *PacketStreamConfig) IsConfigured() bool {
	return i.Configured
}
