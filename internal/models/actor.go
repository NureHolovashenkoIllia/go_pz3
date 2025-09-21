package models

import "time"

// Actor представляє модель актора
type Actor struct {
	ID              uint   `gorm:"primaryKey"`
	FullName        string `gorm:"not null"`
	BirthDate       time.Time
	ExperienceYears int
	ContactInfo     string

	// Зв'язки
	Roles      []Role       // One-to-Many: Один актор може мати багато ролей
	Rehearsals []*Rehearsal `gorm:"many2many:rehearsal_actors;"` // Many-to-Many
}
