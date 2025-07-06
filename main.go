package main

import (
	"fmt"
	"time"

	"github.com/iswangyf/config_center_service/internal/dbinit"
	"github.com/iswangyf/config_center_service/internal/model"
	"github.com/iswangyf/config_center_service/internal/repository"
)

func testInsertAndQueryModuleGroup() {
	dbinit.InitDB()
	if dbinit.DB == nil {
		fmt.Println("Database initialization failed.")
		return
	}
	// Initialize the repository with the database connection
	repo := repository.GetConfigRepository(dbinit.DB)
	group := &model.ModuleGroup{
		ID:          20,
		Name:        "test group",
		Description: "desc",
	}
	if err := repo.InsertModuleGroup(group); err != nil {
		fmt.Printf("InsertModuleGroup failed: %v\n", err)
		return
	}
	got, err := repo.QueryModuleGroupByID(group.ID)
	if err != nil {
		fmt.Printf("QueryModuleGroupByID failed: %v\n", err)
		return
	}
	if got.Name != group.Name {
		fmt.Printf("expected name %s, got %s\n", group.Name, got.Name)
	}
}

func testInsertAndQueryModule() {
	dbinit.InitDB()
	if dbinit.DB == nil {
		fmt.Println("Database initialization failed.")
		return
	}
	// Initialize the repository with the database connection
	repo := repository.GetConfigRepository(dbinit.DB)
	group := &model.ModuleGroup{Name: "g", Description: "d"}
	if err := repo.InsertModuleGroup(group); err != nil {
		fmt.Printf("InsertModuleGroup failed: %v\n", err)
		return
	}
	module := &model.Module{
		GroupID:   group.ID,
		Name:      "mod1",
		Content:   "content",
		ValidFrom: time.Now(),
		ValidTo:   time.Now().AddDate(1, 0, 0), // Valid for one year
		Enabled:   true,
	}
	if err := repo.InsertModule(module); err != nil {
		fmt.Printf("InsertModule failed: %v\n", err)
		return
	}
	modules, err := repo.QueryModulesByGroupID(group.ID)
	if err != nil {
		fmt.Printf("QueryModulesByGroupID failed: %v\n", err)
		return
	}
	if len(modules) == 0 {
		fmt.Println("No modules found for the group.")
	} else {
		fmt.Printf("Found %d modules for group %s.\n", len(modules), group.Name)
	}
}

func main() {
	fmt.Println("Config Center Service started.")
	testInsertAndQueryModuleGroup()
	testInsertAndQueryModule()
	fmt.Println("Test completed.")
}
