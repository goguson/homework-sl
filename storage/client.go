package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"io"
)

const (
	bucketName        = "homework-object-storage"
	expectedNodeCount = 3
)

type NodeClient interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) (err error)
}

func NewNodeClient(ctx context.Context, s Socket, objectID string) (NodeClient, error) {
	details, err := s.FetchNodeDescriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("FetchNodeDescriptions: %w", err)
	}

	d, err := details.SelectDescriptionByID(objectID)
	log.Ctx(ctx).Info().Msgf("Selected node: %s", d.ip)
	if err != nil {
		return nil, fmt.Errorf("SelectDescriptionByID: %w", err)
	}
	return minio.New(d.ip, &minio.Options{
		Creds: credentials.NewStaticV4(d.accessKey, d.secretKey, ""),
	})
}
