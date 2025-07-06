package repository

import (
	"testing"
	"time"

	"github.com/iswangyf/config_center_service/internal/model"
	"gorm.io/driver/sqlite" // Required CGO_ENABLED=1 for SQLite, and GCC
	"gorm.io/gorm"
)

func setupTestRepo(t *testing.T) *ConfigRepository {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	// Auto migrate the models
	if err := db.AutoMigrate(&model.ModuleGroup{}, &model.Module{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	// Clear the database before each test
	if err := db.Exec("DELETE FROM module_groups").Error; err != nil {
		t.Fatalf("failed to clear module_groups: %v", err)
	}
	if err := db.Exec("DELETE FROM modules").Error; err != nil {
		t.Fatalf("failed to clear modules: %v", err)
	}
	// Return a new ConfigRepository instance
	return GetConfigRepository(db)
}

func TestInsertAndQueryModuleGroup(t *testing.T) {
	repo := setupTestRepo(t)
	group := &model.ModuleGroup{
		Name:        "test group",
		Description: "desc",
		CreatedAt:   time.Now(),
	}
	err := repo.InsertModuleGroup(group)
	if err != nil {
		t.Fatalf("InsertModuleGroup failed: %v", err)
	}
	got, err := repo.QueryModuleGroupByID(group.ID)
	if err != nil {
		t.Fatalf("QueryModuleGroupByID failed: %v", err)
	}
	if got.Name != group.Name {
		t.Errorf("expected name %s, got %s", group.Name, got.Name)
	}
}

func TestInsertAndQueryModule(t *testing.T) {
	repo := setupTestRepo(t)
	group := &model.ModuleGroup{Name: "g", Description: "d", CreatedAt: time.Now()}
	if err := repo.InsertModuleGroup(group); err != nil {
		t.Fatalf("InsertModuleGroup failed: %v", err)
	}
	module := &model.Module{
		GroupID:   group.ID,
		Name:      "mod1",
		Content:   "content",
		ValidFrom: time.Now(),
		ValidTo:   time.Now().Add(24 * time.Hour),
		Enabled:   true,
		CreatedAt: time.Now(),
	}
	if err := repo.InsertModule(module); err != nil {
		t.Fatalf("InsertModule failed: %v", err)
	}
	list, err := repo.QueryModulesByGroupID(group.ID)
	if err != nil {
		t.Fatalf("QueryModulesByGroupID failed: %v", err)
	}
	if len(list) != 1 || list[0].Name != "mod1" {
		t.Errorf("expected 1 module named mod1, got %+v", list)
	}
}

func TestUpdateAndDeleteModule(t *testing.T) {
	repo := setupTestRepo(t)
	group := &model.ModuleGroup{Name: "g", Description: "d", CreatedAt: time.Now()}
	_ = repo.InsertModuleGroup(group)
	module := &model.Module{
		GroupID:   group.ID,
		Name:      "mod1",
		Content:   "content",
		ValidFrom: time.Now(),
		ValidTo:   time.Now().Add(24 * time.Hour),
		Enabled:   true,
		CreatedAt: time.Now(),
	}
	_ = repo.InsertModule(module)
	module.Name = "mod2"
	if err := repo.UpdateModule(module); err != nil {
		t.Fatalf("UpdateModule failed: %v", err)
	}
	got, err := repo.QueryModulesByGroupID(group.ID)
	if err != nil || got[0].Name != "mod2" {
		t.Errorf("UpdateModule did not update name, got %+v", got)
	}
	if err := repo.DeleteModule(module.ID); err != nil {
		t.Fatalf("DeleteModule failed: %v", err)
	}
	got, err = repo.QueryModulesByGroupID(group.ID)
	if err != nil {
		t.Fatalf("QueryModulesByGroupID failed: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 modules after delete, got %+v", got)
	}
}
