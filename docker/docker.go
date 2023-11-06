package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
)

func GetDockerClient() (*client.Client, error) {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %v", err)
	}

	_, err = dockerClient.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect Docker server: %v", err)
	}

	return dockerClient, nil
}
