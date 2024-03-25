package main

import (
	"errors"

	"github.com/atotto/clipboard"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
)

const (
	EARNAPP_IMAGE_NAME = "fazalfarhan01/earnapp:lite"
)

type EarnAppConfig struct {
	UUID       string
	Configured bool
}

func (i *EarnAppConfig) ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application) {
	uuid := ""
	isError := false
	showingError := false
	form.AddInputField("UUID", i.UUID, 50, nil, func(text string) {
		uuid = text
	})
	form.AddTextView("Attention", "The Claim URL will work only after the container has been started,\nso copy it and keep it in a safe place until you need it", 0, 0, true, false)
	form.AddButton("Generate UUID", func() {
		uuid = generateEarnAppUUID()
		form.GetFormItemByLabel("UUID").(*tview.InputField).SetText(uuid)
	})
	form.AddButton("Copy Claim URL to Clipboard", func() {
		if stringIsEmpty(uuid) {
			form.AddTextView("Error", "UUID is required", 0, 1, true, true)
			isError = true
		} else {
			if err := clipboard.WriteAll("https://earnapp.com/r/" + uuid); err != nil {
				form.AddTextView("Error", "Failed to copy to clipboard", 0, 1, true, true)
			} else {
				form.AddTextView("Success", "Copied to clipboard", 0, 1, true, true)
			}
		}
	})
	form.AddButton("Save", func() {
		isError = stringIsEmpty(uuid)
		if isError {
			if !showingError {
				form.AddTextView("Error", "All fields are required", 0, 1, true, true)
				showingError = true
			}
			return
		}
		i.UUID = uuid
		i.Configured = true
		returnToMenu(frame, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(frame, app)
	})
	form.AddButton("Register", func() {
		webbrowser.Open("https://earnapp.com/i/J9XF4PXJ")
	})
}

func (i *EarnAppConfig) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `earnapp:
  image: ` + EARNAPP_IMAGE_NAME + `
  environment:
    - EARNAPP_UUID=` + i.UUID + `
    - EARNAPP_TERM="yes"
  volumes:
    - earnapp-data:/etc/earnapp
  restart: unless-stopped
`, nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		containerConfig := &container.Config{
			Image: EARNAPP_IMAGE_NAME,
			Env: []string{
				"EARNAPP_UUID=" + i.UUID,
				"EARNAPP_TERM=yes",
			},
		}
		hostConfig := &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: "earnapp-data",
					Target: "/etc/earnapp",
				},
			},
		}
		return "", createContainer("earnapp", containerConfig, hostConfig, logView)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *EarnAppConfig) IsConfigured() bool {
	return i.Configured
}
