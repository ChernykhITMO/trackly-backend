package repositories

import (
	"gorm.io/gorm"
	"trackly-backend/internal/models"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {

	return &PlanRepository{db: db}
}

func (p *PlanRepository) CreatePlan(plan *models.Plan) error {
	return p.db.Create(plan).Error
}

func (p *PlanRepository) GetPlansByHabitId(id int) (*[]models.Plan, error) {

	var plan []models.Plan

	if err := p.db.Where("habit_id = ?", id).Find(&plan).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (p *PlanRepository) UpdatePlan(plan *models.Plan) error {

	if err := p.db.Save(plan).Error; err != nil {
		return err
	}
	return nil
}
