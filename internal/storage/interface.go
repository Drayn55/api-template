package storage

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	Upload(
		ctx context.Context,
		file multipart.File,
		fileName string,
		contentType string,
	) (string, error)

	Delete(
		ctx context.Context,
		objectKey string,
	) error
}
