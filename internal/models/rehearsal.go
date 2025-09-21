package models

import "time"

// Rehearsal представляє модель репетиції
type Rehearsal struct {
	ID       uint      `gorm:"primaryKey"`
	DateTime time.Time `gorm:"not null"`

	// Зовнішній ключ для зв'язку "належить до" (Belongs To)
	PerformanceID uint

	// Зв'язки
	Actors []*Actor `gorm:"many2many:rehearsal_actors;"` // Many-to-Many
}
