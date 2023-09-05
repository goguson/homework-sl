package storage

import (
	"cmp"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
	"hash/fnv"
	"slices"
	"strings"
)

const (
	imageName      = "minio/minio"
	network        = "homework-object-storage_amazin-object-storage"
	defaultApiPort = ":9000"
)

type Fetcher interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}

// Socket most likely does not need to be thread safe
type Socket struct {
	d Fetcher
}

type NodeDescription struct {
	id        string
	ip        string
	accessKey string
	secretKey string
}

type NodeDescriptions map[int]NodeDescription

func NewSocketClient(host string) (Socket, error) {
	d, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return Socket{}, fmt.Errorf("NewClientWithOpts: %w", err)
	}
	return Socket{d: d}, nil

}
func (nodes NodeDescriptions) SelectDescriptionByID(id string) (NodeDescription, error) {
	k, err := pickNodeNumber(id, expectedNodeCount)
	if err != nil {
		return NodeDescription{}, fmt.Errorf("pickNodeNumber: %w", err)
	}

	node, ok := nodes[k]
	if !ok {
		return NodeDescription{}, fmt.Errorf("instance not accessible")
	}

	return node, nil
}

func (s *Socket) FetchNodeDescriptions(ctx context.Context) (NodeDescriptions, error) {
	res, err := s.d.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("ContainerList: %w", err)
	}

	var results []NodeDescription
	for _, re := range res {
		if re.Image != imageName {
			continue
		}

		r, err := s.d.ContainerInspect(ctx, re.ID)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("ContainerInspect: %s", err.Error())
			continue
		}

		net, ok := r.NetworkSettings.Networks[network]
		if !ok {
			log.Ctx(ctx).Error().Msgf("ContainerInspect: %s network not found", network)
			continue
		}

		var secretKey string
		var accessKey string
		for _, e := range r.Config.Env {
			if secretKey != "" && accessKey != "" {
				break
			}
			if strings.Contains(e, "MINIO_SECRET_KEY=") {
				secretKey = e[17:]
				continue
			}
			if strings.Contains(e, "MINIO_ACCESS_KEY=") {
				accessKey = e[17:]
				continue
			}
		}

		i, err := newNode(r.ID, net.IPAddress, accessKey, secretKey)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("newNode: %s", err.Error())
			continue
		}
		results = append(results, i)
	}

	slices.SortFunc(results, func(left, right NodeDescription) int {
		return cmp.Compare(strings.ToLower(left.accessKey), strings.ToLower(right.accessKey))
	})
	m := map[int]NodeDescription{}
	for i, result := range results {
		m[i] = result
	}

	return m, nil
}

func newNode(id, ip, accessKey, secretKey string) (NodeDescription, error) {
	if id == "" || ip == "" || accessKey == "" || secretKey == "" {
		return NodeDescription{},
			fmt.Errorf("at least one missing field value: id: %s, ip: %s, accessKey: %s, secretKey: %s",
				id, ip, accessKey, secretKey)
	}
	return NodeDescription{
		id:        id,
		ip:        ip + defaultApiPort,
		accessKey: accessKey,
		secretKey: secretKey,
	}, nil
}

func pickNodeNumber(objectID string, instanceCount int) (int, error) {
	if objectID == "" {
		return 0, fmt.Errorf("objectID is empty")
	}

	if instanceCount < 2 {
		return 0, fmt.Errorf("instanceCount is below minimal value: 2")
	}

	hash := fnv.New32a()
	_, err := hash.Write([]byte(objectID))
	if err != nil {
		return 0, fmt.Errorf("CalculateHash: %w", err)
	}
	hashValue := int(hash.Sum32())
	return hashValue % instanceCount, nil
}
