package repository

import (
	"github.com/iswangyf/config_center_service/internal/model"
	"gorm.io/gorm"
)

// ConfigRepository provides methods to interact with the configuration database.
type ConfigRepository struct {
	db *gorm.DB
}

// GetConfigRepository returns a singleton instance of ConfigRepository.
func GetConfigRepository(configDb *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db: configDb}
}

// QueryModuleGroups retrieves all module groups from the database.
func (r *ConfigRepository) QueryModuleGroups() ([]model.ModuleGroup, error) {
	var groups []model.ModuleGroup
	if err := r.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// QueryModuleGroupByID retrieves a module group by its ID.
func (r *ConfigRepository) QueryModuleGroupByID(id uint) (*model.ModuleGroup, error) {
	var group model.ModuleGroup
	if err := r.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// QueryModulesByGroupID retrieves all modules for a given module group ID.
func (r *ConfigRepository) QueryModulesByGroupID(groupID uint) ([]model.Module, error) {
	var modules []model.Module
	if err := r.db.Where("group_id = ?", groupID).Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

// InsertModuleGroup inserts a new module group into the database.
func (r *ConfigRepository) InsertModuleGroup(group *model.ModuleGroup) error {
	if err := r.db.Create(group).Error; err != nil {
		return err
	}
	return nil
}

// InsertModule inserts a new module into the database.
func (r *ConfigRepository) InsertModule(module *model.Module) error {
	if err := r.db.Create(module).Error; err != nil {
		return err
	}
	return nil
}

// UpdateModule updates an existing module in the database.
func (r *ConfigRepository) UpdateModule(module *model.Module) error {
	if err := r.db.Save(module).Error; err != nil {
		return err
	}
	return nil
}

// DeleteModule deletes a module by its ID from the database.
func (r *ConfigRepository) DeleteModule(id uint) error {
	if err := r.db.Delete(&model.Module{}, id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteModuleGroup deletes a module group by its ID from the database.
func (r *ConfigRepository) DeleteModuleGroup(id uint) error {
	if err := r.db.Delete(&model.ModuleGroup{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Close closes the database connection.
func (r *ConfigRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the underlying gorm.DB instance.
func (r *ConfigRepository) GetDB() *gorm.DB {
	if r.db == nil {
		panic("Database connection is not initialized")
	}
	return r.db
}
