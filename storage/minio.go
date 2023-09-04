package storage

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"strings"
)

const imageName = "minio/minio"
const network = "homework-object-storage_amazin-object-storage"

type DockerFetcher interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}

type Client struct {
	c DockerFetcher
}

type Instance struct {
	ID        string
	IP        string
	accessKey string
	secretKey string
}

func NewClient(c DockerFetcher) *Client {
	return &Client{c: c}
}

func NewInstance(id, ip, accessKey, secretKey string) (Instance, error) {
	if id == "" || ip == "" || accessKey == "" || secretKey == "" {
		return Instance{},
			fmt.Errorf("at least one missing field value: id: %s, ip: %s, accessKey: %s, secretKey: %s",
				id, ip, accessKey, secretKey)
	}
	return Instance{
		ID:        id,
		IP:        ip,
		accessKey: accessKey,
		secretKey: secretKey,
	}, nil
}

// FetchMinioInstances make use of network and image name to filter out the instances
func (s *Client) FetchMinioInstances(ctx context.Context) ([]Instance, error) {
	res, err := s.c.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("ContainerList: %w", err)
	}
	var results []Instance
	for _, re := range res {
		if re.Image != imageName {
			continue
		}

		r, err := s.c.ContainerInspect(ctx, re.ID)
		if err != nil {
			continue
		}

		net, ok := r.NetworkSettings.Networks[network]
		if !ok {
			continue
		}

		var secretKey string
		var accessKey string
		for _, e := range r.Config.Env {
			if strings.Contains(e, "MINIO_SECRET_KEY=") {
				secretKey = e[17:]
			}
			if strings.Contains(e, "MINIO_ACCESS_KEY=") {
				accessKey = e[17:]
			}
		}
		i, err := NewInstance(r.ID, net.IPAddress, accessKey, secretKey)
		if err != nil {
			// we can have a log for tracing semi-broken instances for future investigation
			continue
		}
		results = append(results, i)
	}
	return results, nil
}
