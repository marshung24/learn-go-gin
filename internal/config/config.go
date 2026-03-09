package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 包含應用程式所有設定
type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
}

// AppConfig 應用程式設定
type AppConfig struct {
	Name string
	Port string
}

// DBConfig 資料庫連線設定
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// RedisConfig Redis 連線設定
type RedisConfig struct {
	Host string
	Port string
}

// AppConfig 全域設定實例
var AppCfg Config

// LoadConfig 載入設定檔
// Viper 支援多種來源：.env 檔案、環境變數、YAML 等
// AutomaticEnv() 讓環境變數自動覆寫檔案中的設定
func LoadConfig() {
	// 設定 .env 檔案位置
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// 環境變數自動綁定
	// 若環境變數存在，會覆蓋 .env 檔案中的值
	viper.AutomaticEnv()

	// 設定預設值
	setDefaults()

	// 嘗試讀取 .env 檔案
	// 讀取失敗不算錯誤（可能只用環境變數）
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	}

	// 綁定到 Config struct
	AppCfg = Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Port: viper.GetString("APP_PORT"),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		Redis: RedisConfig{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
	}

	log.Printf("Config loaded: App=%s, Port=%s, DB=%s:%s/%s, Redis=%s:%s\n",
		AppCfg.App.Name, AppCfg.App.Port,
		AppCfg.DB.Host, AppCfg.DB.Port, AppCfg.DB.Name,
		AppCfg.Redis.Host, AppCfg.Redis.Port)
}

// setDefaults 設定預設值
func setDefaults() {
	viper.SetDefault("APP_NAME", "go-gin-app")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "app")
	viper.SetDefault("DB_PASSWORD", "app")
	viper.SetDefault("DB_NAME", "app")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
}

// GetDSN 取得 MySQL Data Source Name（連線字串）
func GetDSN() string {
	return AppCfg.DB.User + ":" + AppCfg.DB.Password +
		"@tcp(" + AppCfg.DB.Host + ":" + AppCfg.DB.Port + ")/" +
		AppCfg.DB.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// GetRedisAddr 取得 Redis 連線位址
func GetRedisAddr() string {
	return AppCfg.Redis.Host + ":" + AppCfg.Redis.Port
}
