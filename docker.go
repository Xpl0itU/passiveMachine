package main

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/rivo/tview"
)

type DockerConfigKind string

const (
	KIND_DOCKER_COMPOSE            = "docker-compose"
	KIND_DIRECTLY_CONFIGURE_DOCKER = "directly-configure-docker"
)

var dockerClient *client.Client

func getDockerClient() (*client.Client, error) {
	if dockerClient == nil {
		var err error
		dockerClient, err = client.NewClientWithOpts(client.WithHostFromEnv(), client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
	}
	return dockerClient, nil
}

func buildDockerComposeFile(menuItems []MenuItem) string {
	dockerComposeFile := `version: '3'
services:
`
	for _, item := range menuItems {
		if !item.Config.IsConfigured() {
			continue
		}
		dockerCompose, err := item.Config.ConfigureDocker(KIND_DOCKER_COMPOSE, nil)
		if err != nil {
			continue
		}
		dockerComposeFile += dockerCompose
	}
	return dockerComposeFile
}

func batchCreateDockerContainers(menuItems []MenuItem, logView *tview.TextView) []error {
	var errors []error
	for _, item := range menuItems {
		if !item.Config.IsConfigured() {
			continue
		}
		_, err := item.Config.ConfigureDocker(KIND_DIRECTLY_CONFIGURE_DOCKER, logView)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func pullImageBlocking(imageName string, logView *tview.TextView) error {
	client, err := getDockerClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	out, err := client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	// log the output
	logView.Clear()
	logView.ScrollToEnd()
	logView.Write([]byte("Pulling image " + imageName + "\n"))
	logView.ScrollToEnd()
	buf := make([]byte, 1024)
	for {
		n, err := out.Read(buf)
		if n > 0 {
			logView.Write(buf[:n])
			logView.ScrollToEnd()
		}
		if err != nil {
			break
		}
	}

	out.Close()

	return nil
}

func createContainer(name string, containerConfig *container.Config, hostConfig *container.HostConfig, logView *tview.TextView) error {
	client, err := getDockerClient()
	if err != nil {
		return err
	}

	if err := pullImageBlocking(containerConfig.Image, logView); err != nil {
		return err
	}

	out, err := client.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, name)
	if err != nil {
		return err
	}
	for _, warning := range out.Warnings {
		logView.Write([]byte(warning))
		logView.ScrollToEnd()
	}
	return client.ContainerStart(context.Background(), out.ID, container.StartOptions{})
}
