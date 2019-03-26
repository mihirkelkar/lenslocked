package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/* This services package was created to instantiate all our model
service the moment we get data base credentials. As we grow, we can't
keep instantiating more database connections everytime we
add a new service. So this service starts one single database connection
and then instantiates all models here
*/

//Services : Single enclosure to instantiate and store all kinds of
//resource services
type Services struct {
	UserService    UserService
	GalleryService GalleryService
	ImageService   ImageService
	db             *gorm.DB
}

//NewServices : Instatntiates a bunch of different services here.
// and then returns an enclosing struct containing them.
//we're not returning an interface here, we're returning a pointer to a struct
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	userService, errUser := NewUserService(db)
	if errUser != nil {
		return nil, errUser
	}
	galleryService, errGall := NewGalleryService(db)
	if errGall != nil {
		return nil, errGall
	}

	imageService, errImg := NewImageService()
	if errImg != nil {
		return nil, errImg
	}

	return &Services{
		UserService:    userService,
		GalleryService: galleryService,
		ImageService:   imageService,
		db:             db,
	}, nil
}

//Close : Closes the database connection, in turn closing all services
func (s *Services) Close() error {
	return s.db.Close()
}

//AutoMigrate : Migrate the models in the database. This adds / removes columns
// from the database models when we make updates.
func (s *Services) AutoMigrate() error {
	err := s.db.AutoMigrate(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return nil
}

//Destructive Reset : Destroy the whole database.
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	s.AutoMigrate()
	return nil
}
