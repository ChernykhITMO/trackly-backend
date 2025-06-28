package models

import "time"

type Habit struct {
	ID            int          `gorm:"primaryKey" json:"habit_id"`
	HabitName     string       `json:"habit_name"`
	Description   *string      `json:"description"`
	UserId        int          `json:"user_id"`
	StartDate     time.Time    `json:"start_date"`
	Notifications *bool        `json:"notifications"`
	Plans         []Plan       `gorm:"foreignKey:HabitId"`
	HabitScore    []HabitScore `gorm:"foreignKey:HabitId"`
}
