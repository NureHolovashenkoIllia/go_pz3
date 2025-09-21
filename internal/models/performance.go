package models

import "time"

// Performance представляє модель вистави
type Performance struct {
	ID              uint   `gorm:"primaryKey"`
	Title           string `gorm:"not null"`
	PremiereDate    time.Time
	Genre           string
	DurationMinutes int

	// Зв'язки
	Roles      []Role      // One-to-Many: Одна вистава має багато ролей
	Rehearsals []Rehearsal // One-to-Many: Одна вистава має багато репетицій
}
