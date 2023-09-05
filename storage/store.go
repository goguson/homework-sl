package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type Store interface {
	Put(ctx context.Context, objectID string, r *http.Request) error
	Get(ctx context.Context, objectID string, w http.ResponseWriter, r *http.Request) error
}

type storage struct {
	s Socket
}

func NewStore(s Socket) Store {
	return &storage{s: s}
}

func (c storage) Put(ctx context.Context, objectID string, r *http.Request) error {
	node, err := NewNodeClient(ctx, c.s, objectID)
	if err != nil {
		return fmt.Errorf("NewNodeClient: %w", err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("FormFile: %w", err)
	}
	defer file.Close()

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

func (c storage) Get(ctx context.Context, objectID string, w http.ResponseWriter, r *http.Request) error {

	node, err := NewNodeClient(ctx, c.s, objectID)
	if err != nil {
		return fmt.Errorf("NewNodeClient: %w", err)
	}

	reader, err := node.GetObject(ctx, bucketName, objectID, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("GetObject: %w", err)
	}

	defer reader.Close()

	stat, err := reader.Stat()
	if err != nil {
		log.Ctx(ctx).Err(fmt.Errorf("reader.Stat: %w", err)).Send()
		return minio.ToErrorResponse(err)
	}

	w.Header().Set("Content-Type", stat.ContentType)

	_, err = io.Copy(w, reader)
	if err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	return nil
}
