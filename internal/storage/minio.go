package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioStorage(
	client *minio.Client,
	bucketName string,
) *MinioStorage {
	return &MinioStorage{
		Client:     client,
		BucketName: bucketName,
	}
}

func (m *MinioStorage) Upload(
	ctx context.Context,
	file multipart.File,
	fileName string,
	contentType string,
) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("minio client not configured")
	}

	objectKey := fileName

	_, err := m.Client.PutObject(
		ctx,
		m.BucketName,
		objectKey,
		file,
		-1,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", err
	}

	return objectKey, nil
}

func (m *MinioStorage) Delete(
	ctx context.Context,
	objectKey string,
) error {
	if m.Client == nil {
		return fmt.Errorf("minio client not configured")
	}

	return m.Client.RemoveObject(
		ctx,
		m.BucketName,
		objectKey,
		minio.RemoveObjectOptions{},
	)
}

func (m *MinioStorage) GetPublicURL(
	objectKey string,
	endpoint string,
) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		endpoint,
		m.BucketName,
		objectKey,
	)
}

func (m *MinioStorage) GeneratePresignedUploadURL(
	ctx context.Context,
	objectKey string,
	expiry time.Duration,
) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("minio client not configured")
	}

	if expiry <= 0 {
		expiry = 15 * time.Minute
	}

	u, err := m.Client.PresignedPutObject(
		ctx,
		m.BucketName,
		objectKey,
		expiry,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload url: %w", err)
	}

	return u.String(), nil
}

