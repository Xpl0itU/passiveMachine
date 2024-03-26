package main

import (
	"errors"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/rivo/tview"
	"github.com/toqueteos/webbrowser"
)

const (
	PEER2PROFIT_IMAGE_NAME = "enwaiax/peer2profit:latest"
)

type Peer2ProfitConfig struct {
	Email      string
	Configured bool
}

func (i *Peer2ProfitConfig) ConfigureForm(form *tview.Form, frame *tview.Frame, app *tview.Application) {
	email := ""
	isError := false
	showingError := false
	form.AddInputField("Email", i.Email, 50, nil, func(text string) {
		email = text
	})
	form.AddButton("Save", func() {
		isError = stringIsEmpty(email)
		if isError {
			if !showingError {
				form.AddTextView("Error", "All fields are required", 0, 1, true, true)
				showingError = true
			}
			return
		}
		i.Email = email
		i.Configured = true
		returnToMenu(frame, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(frame, app)
	})
	form.AddButton("Register", func() {
		webbrowser.Open("https://t.me/peer2profit_app_bot?start=1671204644639c8f24d663c")
	})
}

func (i *Peer2ProfitConfig) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		compose := `peer2profit:
  image: ` + PEER2PROFIT_IMAGE_NAME + `
  environment:
	- email=` + i.Email + `
    - use_proxy=false
  restart: unless-stopped
  platform: linux/amd64
`
		if runtime.GOARCH == "arm64" {
			return compose + `binfmt:
  image: tonistiigi/binfmt:latest
  privileged: true
  command: --install all
  restart: unless-stopped
`, nil
		}
		return compose, nil

	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		if runtime.GOARCH == "arm64" {
			containerConfig := &container.Config{
				Image: "tonistiigi/binfmt:latest",
				Env:   []string{},
				Cmd:   []string{"--install", "all"},
			}
			hostConfig := &container.HostConfig{
				Privileged: true,
				RestartPolicy: container.RestartPolicy{
					Name: "unless-stopped",
				},
			}
			if err := createContainer("binfmt", containerConfig, hostConfig, logView); err != nil {
				return "", err
			}
		}
		containerConfig := &container.Config{
			Image: PEER2PROFIT_IMAGE_NAME,
			Env: []string{
				"email=" + i.Email,
				"use_proxy=false",
			},
		}
		hostConfig := &container.HostConfig{
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		}
		return "", createContainer("peer2profit", containerConfig, hostConfig, logView)
	}
	return "", errors.New("unknown kind")
}

func (i *Peer2ProfitConfig) IsConfigured() bool {
	return i.Configured
}
