package repositories

import (
	"gorm.io/gorm"
	"trackly-backend/internal/models"
)

type StatisticRepository struct {
	db *gorm.DB
}

func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{db: db}
}

func (s *StatisticRepository) CreateStatistic(score *models.HabitScore) error {
	return s.db.Create(score).Error
}

func (s *StatisticRepository) GetAllStatisticByHabitId(habit int) ([]models.Plan, error) {
	var plan []models.Plan

	if err := s.db.Where("habit_id = ?", habit).Find(&plan).Error; err != nil {
		return nil, err
	}
	return plan, nil
}
