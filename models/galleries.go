package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not_null;index"`
	Title  string  `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

var (
	ErrUserIDRequired modelerror = "Error: UserID is reqquired"
	ErrTitleRequired  modelerror = "Error: A title is required"
	ErrIdNotFound     modelerror = "Error: The Gallery ID was not found"
	ErrZeroID         modelerror = "Error: Zero Gallery ID cannot be deleted"
)

type GalleryValFns func(*Gallery) error

func RunGalleryValFns(gallery *Gallery, fns ...GalleryValFns) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
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
	ByID(id uint) (*Gallery, error)
	ByUserID(id uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(gallery *Gallery) error
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

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(gallery *Gallery) error {
	return gg.db.Delete(gallery).Error
}

//ByID : Reciever function defined on galleryGorm that fits the GalleryDB interface
func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	err := gg.db.Where("id = ?", id).First(&gallery).Error
	switch err {
	case nil:
		return &gallery, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrIdNotFound
	default:
		return nil, ErrIdNotFound
	}
}

//ByUserID : Finds all galleries associated with a userId
func (gg *galleryGorm) ByUserID(id uint) ([]Gallery, error) {
	var galleries []Gallery
	err := gg.db.Where("user_id= ?", id).Find(&galleries).Error
	switch err {
	case nil:
		return galleries, err
	case gorm.ErrRecordNotFound:
		return nil, ErrIdNotFound
	default:
		return nil, ErrIdNotFound
	}
}

//NewGalleryService : Retrurns a new gallery service letting controllers use this to
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
func (gv *galleryValidator) userIDRequired(gallery *Gallery) error {
	//check if the userId of the gallery is nil. If so return error.
	//otherwise return nil.
	//In gorm, we start user ids from 1
	if gallery.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

//This function is a reciever funcion on the gallery validator struct
//It also implements the GalleryValFns type.
func (gv *galleryValidator) titleRequried(gallery *Gallery) error {
	if gallery.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) nonZeroID(gallery *Gallery) error {
	if gallery.ID == 0 {
		return ErrZeroID
	}
	return nil
}

//This runs all the validation function and just calls the underlying create
func (gv *galleryValidator) Create(gallery *Gallery) error {
	if err := RunGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequried); err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

//This runs all the validation functions and just calls the underlying update
func (gv *galleryValidator) Update(gallery *Gallery) error {
	if err := RunGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequried); err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

func (gv *galleryValidator) Delete(gallery *Gallery) error {
	if err := RunGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequried,
		gv.nonZeroID); err != nil {
		return err
	}
	return gv.GalleryDB.Delete(gallery)
}

//adds the images to several differnt columns to be displayed
func (g *Gallery) ImagesSplitN(n int) [][]Image {
	ret := make([][]Image, n)
	//make inner slices of length 0
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}

	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	//return a list of list of images
	return ret
}
