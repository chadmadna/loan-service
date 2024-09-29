package upload

import (
	"fmt"
	"io"
	"os"
	"path"
)

// TODO: Implement actual cloud upload
type UploadService interface {
	UploadFile(file io.Reader, filename, contentType string) (string, error)
}

func NewUploadService() UploadService {
	return &uploadService{}
}

type uploadService struct{}

func (u *uploadService) UploadFile(sourceFile io.Reader, filename, contentType string) (string, error) {
	fmt.Printf("[upload started] %s (%s)\n", filename, contentType)

	// Create new file in public dir
	destinationFile, err := os.Create(path.Join("public", filename))
	if err != nil {
		fmt.Println(err)
	}
	defer destinationFile.Close()

	// Copy source file to destination
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("[upload completed] %s (%s)\n", filename, contentType)

	return destinationFile.Name(), nil
}
