package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

var (
	ErrUserIDRequired modelerror = "Error: UserID is reqquired"
	ErrTitleRequired  modelerror = "Error: A title is required"
)

type GalleryValFns func(*Gallery) error

func RunGalleryValFns(gallery *Gallery, fns ...GalleryValFns) error {
	for _, fn := range(fns){
		if err := fn(gallery); err != nil{
			return err
		}
	}
	return nil
}

//GalleryService : Can be used from the controllers to access the Gallery Gorm Model
type GalleryService interface {
	GalleryDB
}

//Both of these structs below implement the GalleryService Interface.
//They also implement the GAlleryDB interface per se right now since
//they have galleryDB as a field in the struct.
type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
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
	return gg.db.Create(gallery).Error
}

//NewGAlleryService : Retrurns a new gallery service letting controllers use this to
//do things to / with galleries
func NewGalleryService(db *gorm.DB) (GalleryService, error) {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}, nil
}

//This function is a reciever on the gallery validator struct
//and also implements the GalleryValFns type.
func (gv *galleryValidator) userIDRequired(gallery *Gallery) error{
	//check if the userId of the gallery is nil. If so return error.
	//otherwise return nil.
	//In gorm, we start user ids from 1
	if gallery.UserID <= 0{
		return ErrUserIDRequired
	}
	return nil
}

//This function is a reciever funcion on the gallery validator struct
//It also implements the GalleryValFns type.
func (gv *galleryValidator) titleRequried(gallery *Gallery) error{
	if gallery.Title == ""{
		return ErrTitleRequired
	}
	return nil
}

//This runs all the validation function and just calls the underlying create
func (gv *galleryValidator) Create(gallery *Gallery) error {
	if err := RunGalleryValFns(gallery,
		gv.userIDRequired,
	  gv.titleRequried,); err != nil{
			return err
		}
		return gv.GalleryDB.Create(gallery)
}
