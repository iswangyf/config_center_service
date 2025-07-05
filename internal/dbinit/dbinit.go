package dbinit

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/iswangyf/config_center_service/internal/model"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"` // Database name，e.g., configdb, need to be created before running the service
	} `yaml:"database"`
}

var DB *gorm.DB
var AppConfig Config

func InitConfig() {
	file, err := os.ReadFile("config/config.yaml")
	if err != nil {
		panic("Failed to read config file")
	}
	if err := yaml.Unmarshal(file, &AppConfig); err != nil {
		panic("Failed to parse config file")
	}
}

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.DBName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err = DB.AutoMigrate(&model.ModuleGroup{}, &model.Module{}); err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database connection established and migrations completed successfully.")
}
