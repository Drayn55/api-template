// configs/minio.go
package configs

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinIO creates and returns a connected *minio.Client (returns nil if Endpoint is not configured)
func NewMinIO(cfg MinIOConfig) (*minio.Client, error) {
	if cfg.Endpoint == "" {
		return nil, nil
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	return client, nil
}
