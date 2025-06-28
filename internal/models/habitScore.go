package models

import "time"

type HabitScore struct {
	ID       int       `gorm:"primaryKey" json:"score_id"`
	HabitId  int       `json:"habit_id"`
	PlanId   int       `json:"plan_id"`
	DateTime time.Time `json:"date_time"`
	Value    int       `json:"value"`
}
