package repositories

import (
	"gorm.io/gorm"
	"time"
	"trackly-backend/internal/models"
)

type HabitRepository struct {
	db *gorm.DB
}

func NewHabitRepository(db *gorm.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

func (r *HabitRepository) CreateHabit(habit *models.Habit) error {

	return r.db.Create(&habit).Error
}

func (r *HabitRepository) GetHabitsByUserId(userId int) ([]*models.Habit, error) {
	var habits []*models.Habit
	if err := r.db.Preload("HabitScore").Preload("Plans").Where("user_id = ?", userId).Find(&habits).Error; err != nil {
		return nil, err
	}
	return habits, nil
}

func (r *HabitRepository) GetHabitById(id int, userId int) (*models.Habit, error) {
	var habit models.Habit

	if err := r.db.Preload("HabitScore").Preload("Plans").Where("id = ? AND user_id = ?", id, userId).First(&habit).Error; err != nil {
		return nil, err
	}
	return &habit, nil
}

func (r *HabitRepository) GetHabitWithStatInInterval(id int, userId int, startTime, endTime time.Time) (*models.Habit, error) {
	var habit models.Habit

	if err := r.db.Preload("HabitScore", "date_time BETWEEN ? AND ?", startTime, endTime).
		Preload("Plans").
		Where("id = ? AND user_id = ?", id, userId).
		First(&habit).Error; err != nil {
		return nil, err
	}
	return &habit, nil
}

func (r *HabitRepository) DeleteHabitById(id int) error {
	if err := r.db.Delete(&models.Habit{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *HabitRepository) UpdateHabit(habit *models.Habit) error {

	if err := r.db.Save(&habit).Error; err != nil {
		return err
	}
	return nil
}
