package models

// Role представляє зв'язок між актором та виставою
type Role struct {
	ID       uint   `gorm:"primaryKey"`
	RoleName string `gorm:"not null"`

	// Зовнішні ключі для зв'язку "належить до" (Belongs To)
	ActorID       uint
	PerformanceID uint
}
