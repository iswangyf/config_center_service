package model

import (
	"time"
)

type ModuleGroup struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	CreatedAt   time.Time
}

type Module struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	GroupID   uint `gorm:"not null;index"`
	Name      string
	Content   string
	ValidFrom time.Time
	ValidTo   time.Time
	Enabled   bool
	CreatedAt time.Time
}
