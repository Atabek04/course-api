package models

import (
	"gorm.io/gorm"
	"time"
)

type ModuleInfo struct {
	gorm.Model
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ModuleName     string    `json:"module_name"`
	ModuleDuration int       `json:"module_duration"`
	ExamType       string    `json:"exam_type"`
	Version        string    `json:"version"`
}

func (ModuleInfo) TableName() string {
	return "module_info"
}
