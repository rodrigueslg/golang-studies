package pkg

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDir = "go-challenges"
const configsDir = "configs"
const propertiesFile = "properties.env"

type DbConfig struct {
	User     string
	Password string
	Database string
	Driver   string
}

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDir + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	configsPath := string(rootPath) + "/" + configsDir + "/" + propertiesFile

	err := godotenv.Load(configsPath)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetDbConfig() DbConfig {
	loadEnv()
	return DbConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}
}

func GetMysqlDbConnection(config DbConfig) string {
	return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", config.User, config.Password, config.Database)
}

func GetBearerToken() string {
	loadEnv()
	return os.Getenv("BEARER_TOKEN")
}

func GetDbImplementation() string {
	loadEnv()
	return os.Getenv("DB_IMPL")
}
