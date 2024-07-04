package models

import "time"

type User struct {
	ID             uint   `gorm:"primary_key"`
	PassportNumber string `gorm:"unique;not null"`
	Surname        string
	Name           string
	Patronymic     string
	Address        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
