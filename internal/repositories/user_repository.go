package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"trackly-backend/internal/db"
	"trackly-backend/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindUserById(id int) (*models.User, error) {
	return findUser(r.db, withId(id))
}

func (r *UserRepository) FindUserByEmail(username string) (*models.User, error) {
	return findUser(r.db, withUserName(username))
}

func findUser(tx *gorm.DB, opts ...db.CommonScopeOption) (*models.User, error) {
	user := models.User{}
	err := tx.Scopes(opts...).Last(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func withUserName(userName string) db.CommonScopeOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", userName)
	}
}

func withId(uid int) db.CommonScopeOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", uid)
	}
}

func (r *UserRepository) UpdateUserAvatar(userID int, avatarUUID string) error {
	// Assuming `User` is your GORM model representing the `users` table
	user := &models.User{ID: userID} // Assuming `ID` is the primary key field

	// Update the avatar_id field
	result := r.db.Model(user).Update("avatar_id", avatarUUID)
	if result.Error != nil {
		return fmt.Errorf("failed to update avatar for user %d: %v", userID, result.Error)
	}

	// Check if any rows were affected (optional)
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", user.ID)
	}

	return nil
}