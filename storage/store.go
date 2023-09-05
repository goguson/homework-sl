package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
)

type Store interface {
	Put(ctx context.Context, objectID string, file multipart.File, header *multipart.FileHeader) error
	Get(ctx context.Context, objectID string) (*minio.Object, error)
}

type storage struct {
	s Socket
}

func NewStore(s Socket) Store {
	return &storage{s: s}
}

func (c storage) Put(ctx context.Context, objectID string, file multipart.File, header *multipart.FileHeader) error {
	node, err := NewNodeClient(ctx, c.s, objectID)
	if err != nil {
		return fmt.Errorf("NewNodeClient: %w", err)
	}

	exist, err := node.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("BucketExists: %w", err)
	}

	if !exist {
		err = node.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("MakeBucket: %w", err)
		}
	}

	opts := minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}
	_, err = node.PutObject(ctx, bucketName, objectID, file, header.Size, opts)
	if err != nil {
		return fmt.Errorf("PutObject: %w", err)
	}

	return nil
}

func (c storage) Get(ctx context.Context, objectID string) (*minio.Object, error) {
	node, err := NewNodeClient(ctx, c.s, objectID)
	if err != nil {
		return nil, fmt.Errorf("NewNodeClient: %w", err)
	}

	reader, err := node.GetObject(ctx, bucketName, objectID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetObject: %w", err)
	}

	return reader, nil
}
