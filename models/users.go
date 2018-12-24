package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound  = errors.New("The user you were looking for was not found")
	ErrInvalidID = errors.New("The ID you provided is Invalid")
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

type UserService struct {
	db *gorm.DB
}

//NewUserService : Creates a UserService instane with an open connection
func NewUserService(connectionstring string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionstring)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

//Close : Closes the connection to the gorm database
func (u *UserService) Close() {
	u.db.Close()
}

func (u *UserService) ById(id int) (*User, error) {
	var user User
	err := u.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (u *UserService) ByEmail(email string) (*User, error) {
	var user User
	err := u.db.Where("email = ?", email).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

//AutoMigrate : Auto-Migrates the user table and makes new column additions
// and updates
func (u *UserService) AutoMigrate() error {
	err := u.db.AutoMigrate(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

//DestructiveReset : Completely Delete all data and start again
func (u *UserService) DestructiveReset() error {
	if err := u.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	err := u.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Create(user *User) error {
	err := u.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

//UpdateUser : Updates a given user. Updates all the fiels of the user
//depending on the struct you provide.
func (u *UserService) UpdateUser(user *User) error {
	err := u.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

//DeleteUser : the delete user in gorm works in two ways.
// The first way is that you get an id and you delte the id.
//In the second method, you delete the ID =0 to delete all users.
//We can never let the second case happen. So we should write code to avoid that.
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}