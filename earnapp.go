package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/atotto/clipboard"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/rivo/tview"
)

const (
	EARNAPP_IMAGE_NAME = "fazalfarhan01/earnapp:lite"
)

type EarnAppConfig struct {
	UUID string
}

type EarnAppItem struct {
	Name        string
	Description string
	Config      EarnAppConfig
}

func (i *EarnAppItem) GetName() string {
	return i.Name
}

func (i *EarnAppItem) GetDescription() string {
	return i.Description
}

func (i *EarnAppItem) ConfigureForm(form *tview.Form, list *tview.List, app *tview.Application) {
	uuid := ""
	isError := false
	form.AddInputField("UUID", i.Config.UUID, 50, nil, func(text string) {
		uuid = text
	})
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
		if !isError && stringIsEmpty(uuid) {
			form.AddTextView("Error", "All fields are required", 0, 1, true, true)
			isError = true
			return
		}
		i.Config.UUID = uuid
		returnToMenu(list, app)
	})
	form.AddButton("Cancel", func() {
		returnToMenu(list, app)
	})
}

func (i *EarnAppItem) ConfigureDocker(kind DockerConfigKind, logView *tview.TextView) (string, error) {
	switch kind {
	case KIND_DOCKER_COMPOSE:
		return `earnapp:
  image: ` + EARNAPP_IMAGE_NAME + `
  environment:
    - EARNAPP_UUID=` + i.Config.UUID + `
    - EARNAPP_TERM="yes"
  volumes:
    - earnapp-data:/etc/earnapp
  restart: unless-stopped
`, nil
	case KIND_DIRECTLY_CONFIGURE_DOCKER:
		containerConfig := &container.Config{
			Image: EARNAPP_IMAGE_NAME,
			Env: []string{
				"EARNAPP_UUID=" + i.Config.UUID,
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

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateEarnAppUUID() string {
	return fmt.Sprintf("sdk-node-%x", md5.Sum([]byte(randomString(32))))
}
