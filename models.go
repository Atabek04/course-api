package main

import "gorm.io/gorm"

// ModuleInfo represents the structure of data stored in the database
type ModuleInfo struct {
	ID             uint   `json:"id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ModuleName     string `json:"module_name"`
	ModuleDuration int    `json:"module_duration"`
	ExamType       string `json:"exam_type"`
	Version        string `json:"version"`
}

type DBModel struct {
	DB *gorm.DB
}

func (m *DBModel) CreateModule(moduleInfo *ModuleInfo) error {
	return m.DB.Create(moduleInfo).Error
}

func (m *DBModel) ReadAllModules() ([]ModuleInfo, error) {
	var modules []ModuleInfo
	if err := m.DB.Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

func (m *DBModel) UpdateModule(moduleInfo *ModuleInfo) error {
	return m.DB.Save(moduleInfo).Error
}

func (m *DBModel) RemoveModule(id int) error {
	return m.DB.Delete(&ModuleInfo{}, id).Error
}
