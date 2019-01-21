package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

//GalleryService : Can be used from the controllers to access the Gallery Gorm Model
type GalleryService interface {
	GalleryDB
}

//GalleryDB : interface used to access gallery gorm from within the model.
//You can have a validator implement this interface and fit functionality
//and validation with similar function names
type GalleryDB interface {
	Create(gallery *Gallery) error
}

//galleryGorm : actual struct to access the gallery model.
// We will define reciever functions on this that fit the galleryDB interface
type galleryGorm struct {
	db *gorm.DB
}

//Create : Reciever function defined on galleryGorm that fits the GalleryDB interface
func (gg *galleryGorm) Create(gallery *Gallery) error {
	return nil
}

//Retrurns a new gallery service letting controllers use this to
//do things to / with galleries
func NewGalleryService(db *gorm.DB) (GalleryService, error) {
	return nil, nil
}
