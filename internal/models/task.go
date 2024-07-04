package models

import "time"

type Task struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint `gorm:"not null"`
	TaskName  string
	StartTime time.Time
	EndTime   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
