package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	PathInvalidErr modelerror = "Error: Your filepath might be invalid"
)

//Services to set images to a database.
type ImageService interface {
	Create(galleryID uint, r io.Reader, fileName string) error
	ByGalleryID(galleryID uint) ([]string, error)
}

type imageService struct{}

//newImageService: The image service by itself is just an empty struct.
//So we just return an empty struct.
func NewImageService() (ImageService, error) {
	return &imageService{}, nil
}

//Create:  an image file from the uploaded image.
func (is *imageService) Create(galleryID uint, r io.Reader, fileName string) error {
	imagePath, err := is.makeImagePath(galleryID)
	if err != nil {
		return err
	}
	//create the destination file.
	dst, err := os.Create(filepath.Join(imagePath, fileName))
	if err != nil {
		return err
	}

	defer dst.Close()
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	return nil
}

//ByGalleryID : Return a list of images by galleryID.
func (is *imageService) ByGalleryID(galleryID uint) ([]string, error) {
	path := is.imagePath(galleryID)
	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries", fmt.Sprintf("%v", galleryID))
}

func (is *imageService) makeImagePath(galleryID uint) (string, error) {
	imagePath := is.imagePath(galleryID)
	err := os.MkdirAll(imagePath, 0755)
	if err != nil {
		return "", err
	}
	return imagePath, nil
}
