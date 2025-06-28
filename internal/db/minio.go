package db

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"trackly-backend/internal/config"
)

type MinioClient struct {
	Client *minio.Client
}

func NewMinioClient(cfg *config.Config) (*MinioClient, error) {
	client, err := minio.New(cfg.MinioConfig.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioConfig.MinioRootUser, cfg.MinioConfig.MinioRootPassword, ""),
		Secure: cfg.MinioConfig.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	err = initMinio(context.Background(), cfg, client)
	if err != nil {
		return nil, err
	}
	return &MinioClient{Client: client}, nil
}

func initMinio(ctx context.Context, cfg *config.Config, client *minio.Client) error {
	exists, err := client.BucketExists(ctx, cfg.MinioConfig.BucketName)
	if err != nil {
		return err
	}
	if !exists {
		log.Printf("Бакет %s не найден, создаем новый...", cfg.MinioConfig.BucketName)
		err := client.MakeBucket(ctx, cfg.MinioConfig.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
