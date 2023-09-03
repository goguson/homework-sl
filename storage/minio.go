package storage

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const serviceName = "minio/minio"
const network = "homework-object-storage_amazin-object-storage"

type Socket struct {
	c *client.Client
}

func NewSocket(c *client.Client) *Socket {
	return &Socket{c: c}
}

// FetchIPAddresses fetches Addresses based on network name,
// it is said I am not allowed to use copied ipv4 from docker-compose,
// so I will go with the network name instead.
func (s *Socket) FetchIPAddresses(ctx context.Context) ([]string, error) {
	containers, err := s.c.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("ContainerList: %w", err)
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("no running containers: len(containers): 0")
	}

	var minioAddresses []string
	for _, container := range containers {
		if container.Image == "minio/minio" {
			_, ok := container.NetworkSettings.Networks[network]
			if !ok {
				continue
			}
			minioAddresses = append(minioAddresses, container.NetworkSettings.Networks[network].IPAddress)
		}

	}
	if len(minioAddresses) == 0 {
		return nil, fmt.Errorf("no running minio containers with expected network %s", network)
	}

	return minioAddresses, nil
}

func (s *Socket) FetchSecrets(ctx context.Context) ([]string, error) {
	return nil, nil
}
