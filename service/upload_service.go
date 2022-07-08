package service

import (
	"io"
	"mime/multipart"
	"os"
)

type UploadService interface {
	CreateFile(file multipart.File, fileName string) (string, error)
}

type uploadService struct {
}

func NewUploadService() (UploadService, error) {
	return &uploadService{}, nil
}

func (u uploadService) CreateFile(file multipart.File, fileName string) (string, error) {
	// Create file
	dst, err := os.Create(fileName)
	defer dst.Close()
	if err != nil {
		return "", err
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return "Successfully Uploaded File", nil
}
