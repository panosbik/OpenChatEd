package models

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"OpenChatEd/constants"
)

type User struct {
	gorm.Model
	Username         string  `gorm:"uniqueIndex" json:"username"`
	Email            string  `gorm:"uniqueIndex" json:"email"`
	Password         string  `json:"-"`
	IsActive         bool    `gorm:"default:false" json:"-"`
	ProfileImagePath string  `gorm:"" json:"-"`
	ProfileImage     *string `gorm:"-" json:"profileImage"`
}

// Event Hook runs Automatically before create new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Check for duplicate username or email
	var user User
	if err = tx.First(&user, "username = ? OR email = ?", u.Username, u.Email).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		if u.Username == user.Username {
			err = constants.ErrDuplicateUsername
		} else {
			err = constants.ErrDuplicateEmailUser
		}
		return
	}

	// Generate a hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

// Event Hook runs Automatically before create new user
func (u *User) AfterCreate(tx *gorm.DB) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go u.setIcon(&wg)
	wg.Wait()
	return nil
}

// Running After finding user's object
func (u *User) AfterFind(tx *gorm.DB) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go u.setIcon(&wg)
	wg.Wait()
	return nil
}

// Create the url for user icon
func (u *User) setIcon(wg *sync.WaitGroup) {
	defer wg.Done()
	if u.ProfileImagePath != "" {
		url := constants.StoragePath + u.ProfileImagePath
		u.ProfileImage = &url
	}
}

// Helpful method for find user by UUID with specific columns
func (u *User) First(db *gorm.DB, id uint) error {
	if err := db.Omit("password", "is_active", "profile_image_path", "created", "updated", "deleted").
		First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}
	return nil
}
