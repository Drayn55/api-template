package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"api-template/internal/model"
	"api-template/internal/storage"
	"api-template/pkg/apperr"
)

type StorageService struct {
	storage *storage.MinioStorage
}

func NewStorageService(storage *storage.MinioStorage) *StorageService {
	return &StorageService{storage: storage}
}

func (s *StorageService) Upload(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	filePath string,
) (*string, error) {
	if s.storage == nil {
		return nil, apperr.New(apperr.ErrInternal, "storage service is not configured")
	}
	if fileHeader == nil {
		return nil, nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	key, err := s.storage.Upload(
		ctx,
		file,
		filePath,
		fileHeader.Header.Get("Content-Type"),
	)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (s *StorageService) GeneratePresignedUploadURL(
	ctx context.Context,
	req model.PresignedUploadRequest,
) (*model.PresignedUploadResponse, error) {
	if s.storage == nil {
		return nil, apperr.New(apperr.ErrInternal, "storage service is not configured")
	}

	req.Path = strings.ReplaceAll(req.Path, "|", "/")
	req.Folder = strings.ReplaceAll(req.Folder, "|", "/")
	req.SubFolder = strings.ReplaceAll(req.SubFolder, "|", "/")
	req.Filename = strings.ReplaceAll(strings.ReplaceAll(req.Filename, "|", "_"), "/", "_")

	var objectKey string

	if strings.TrimSpace(req.Path) != "" {
		objectKey = strings.TrimPrefix(strings.TrimSpace(req.Path), "/")
	} else if strings.TrimSpace(req.Filename) != "" {
		folder := strings.Trim(strings.TrimSpace(req.Folder), "/")
		subFolder := strings.TrimSpace(req.SubFolder)
		subFolder = strings.ReplaceAll(subFolder, "/", "-")
		subFolder = strings.ReplaceAll(subFolder, "\\", "-")
		subFolder = strings.ReplaceAll(subFolder, "|", "-")
		subFolder = strings.Trim(subFolder, "-")
		filename := strings.TrimSpace(req.Filename)

		now := time.Now()
		year := strings.TrimSpace(req.Year)
		if year == "" {
			year = now.Format("2006")
		}
		month := strings.TrimSpace(req.Month)
		if month == "" {
			month = now.Format("01")
		}

		parts := []string{}
		if folder != "" {
			parts = append(parts, folder)
		}
		parts = append(parts, year, month)
		if subFolder != "" {
			parts = append(parts, subFolder)
		}
		parts = append(parts, filename)

		objectKey = strings.Join(parts, "/")
	} else {
		return nil, apperr.New(apperr.ErrInvalidParam, "either path or filename must be provided")
	}

	objectKey = path.Clean(objectKey)
	objectKey = strings.TrimPrefix(objectKey, "/")

	expiryMinutes := req.ExpiryMinutes
	if expiryMinutes <= 0 {
		expiryMinutes = 15
	}
	expiry := time.Duration(expiryMinutes) * time.Minute

	uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, objectKey, expiry)
	if err != nil {
		return nil, apperr.New(apperr.ErrInternal, fmt.Sprintf("failed to generate presigned url: %v", err))
	}

	fullPath := fmt.Sprintf("%s/%s", s.storage.BucketName, objectKey)

	return &model.PresignedUploadResponse{
		UploadURL:        uploadURL,
		ObjectKey:        objectKey,
		FullPath:         fullPath,
		Bucket:           s.storage.BucketName,
		ExpiresInSeconds: int(expiry.Seconds()),
	}, nil
}
