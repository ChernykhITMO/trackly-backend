package models

import "time"

type Plan struct {
	ID        int        `gorm:"primaryKey" json:"plan_id"`
	HabitId   int        `json:"habit_id"`
	PlanUnit  *PlanUnit  `json:"plan_unit"`
	Goal      *int       `json:"goal"`
	StartTime time.Time  `json:"start_time"`
	CloseTime *time.Time `json:"close_time"`
}

type PlanUnit string

const (
	Count    PlanUnit = "count"
	Distance PlanUnit = "distance"
	Time     PlanUnit = "time"
)
