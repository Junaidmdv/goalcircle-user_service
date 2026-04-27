package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/repository"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
)

type diskStorage struct {
	basePath string
	baseURL  string
	logger   logger.Logger
}

func NewDiskStorage(config *config.DiscStorageConfig, loger logger.Logger) (repository.FileStorage, error) {

	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage dir: %w", err)
	}
	return &diskStorage{basePath: config.BasePath, baseURL: config.BaseUrl}, nil
}

func (d *diskStorage) UploadFile(ctx context.Context, path string, reader io.Reader, meta *repository.FileMetadata) (string, error) {
	fullPath := filepath.Join(d.basePath, filepath.Clean(path))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		d.logger.Error("failed to create directory", "error", err)
		return "", domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	file, err := os.Create(fullPath)
	if err != nil {
		d.logger.Error("failed to create directory", "error", err)
		return "", domain.NewInternalError("Something went wrong. Please try again later", err)

	}
	defer file.Close()
	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(fullPath)
		d.logger.Error("failed to upload file", "error", err)
		return "", domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	url := fmt.Sprintf("%s/%s", d.baseURL, path)
	return url, nil
}
func (d *diskStorage) DeleteFile(ctx context.Context, path string) error {
	fullPath := filepath.Join(d.basePath, filepath.Clean(path))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		d.logger.Error("failed to delete file", "error", err)
		return domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	return nil
}
func (d *diskStorage) GetURL(ctx context.Context, path string) (string, error) {
	fullPath := filepath.Join(d.basePath, filepath.Clean(path))
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		d.logger.Error("failed to get url of file", "error", err)
		return "", domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	return fmt.Sprintf("%s/%s", d.baseURL, path), nil
}
