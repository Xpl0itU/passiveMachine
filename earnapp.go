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
	EARNAPP_IMAGE_NAME    = "fazalfarhan01/earnapp:lite"
	EARNAPP_REFERRAL_LINK = "https://earnapp.com/i/J9XF4PXJ"
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
	form.AddButton("Claim URL", func() {
		isError = stringIsEmpty(uuid)
		if isError {
			if !showingError {
				form.AddTextView("Error", "UUID is required", 0, 1, true, true)
				showingError = true
			}
			return
		}
		modal := tview.NewModal().
			SetText("Claim URL:\nhttps://earnapp.com/r/" + uuid).
			AddButtons([]string{"Copy", "Close"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Copy" {
					clipboard.WriteAll("https://earnapp.com/r/" + uuid)
				}
				app.SetRoot(form, true)
			})
		app.SetRoot(modal, true)
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
		modal := tview.NewModal().
			SetText("Register on EarnApp\n" + EARNAPP_REFERRAL_LINK).
			AddButtons([]string{"Open", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Open" {
					webbrowser.Open(EARNAPP_REFERRAL_LINK)
				}
				app.SetRoot(form, true)
			})
		app.SetRoot(modal, true)
	})
}

func (i *EarnAppConfig) ConfigureDocker(kind DockerConfigKind, form *tview.Form) (string, error) {
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
		return "", createContainer("earnapp", containerConfig, hostConfig, form)
	default:
		return "", errors.New("unknown kind")
	}
}

func (i *EarnAppConfig) IsConfigured() bool {
	return i.Configured
}

func (i *EarnAppConfig) PostConfigure(form *tview.Form, app *tview.Application) {
	form.AddButton("Open EarnApp Claim URL", func() {
		webbrowser.Open("https://earnapp.com/r/" + i.UUID)
	})
}
