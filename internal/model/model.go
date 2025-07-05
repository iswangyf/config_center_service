package model

import (
	"time"
)

type ModuleGroup struct {
	ID          uint
	Name        string
	Description string
	Modules     []Module
}

type Module struct {
	ID        uint
	GroupID   uint
	Name      string
	Content   string
	ValidFrom time.Time
	ValidTo   time.Time
	Enabled   bool
}
