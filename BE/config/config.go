package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	envPath = "./.env"
)

type Env struct {
	AppName     string
	AppAuthor   string
	AppVersion  string
	AppHost     string
	AppPort     string
	ReleaseMode string
}

type DatabaseConfig struct {
	Type         string `json:"Type"`
	Host         string `json:"Host"`
	Port         string `json:"Port"`
	User         string `json:"User"`
	Password     string `json:"Password"`
	Database     string `json:"Database"`
	MaxOpenConns int    `json:"MaxOpenConns"`
	MaxIdleConns int    `json:"MaxIdleConns"`
	MaxLifetime  int    `json:"MaxLifetime"`
	MaxIdleTime  int    `json:"MaxIdleTime"`
}

type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	Db           int
	PoolSize     int
	MinIdleConns int
}

type LogConfig struct {
	Path       string `json:"Path"`
	Level      string `json:"Level"`
	MaxSizeMB  int    `json:"MaxSizeMB"`
	MaxBackups int    `json:"MaxBackups"`
	MaxAgeDays int    `json:"MaxAgeDays"`
}

type Config struct {
	Timezone                   string         `json:"Timezone"`
	SecretKey                  string         `json:"SecretKey"`
	TokenExpire                int            `json:"TokenExpire"`
	RefreshTokenExpire         int            `json:"RefreshTokenExpire"`
	FormatTime                 string         `json:"FormatTime"`
	FormatDate                 string         `json:"FormatDate"`
	MaxTimeoutGracefulShutdown int            `json:"MaxTimeoutGracefulShutdown"`
	Log                        LogConfig      `json:"Log"`
	CorsAllowOrigins           []string       `json:"CorsAllowOrigins"`
	Database                   DatabaseConfig `json:"Database"`
}

var (
	configMap = map[string]string{
		"dev":   "./config/config_dev.json",
		"prod":  "./config/config_prod.json",
		"uat":   "./config/config_uat.json",
		"local": "./config/config_local.json",
		"bors":  "./config/config_bors_kost.json",
	}
	ENV        *Env
	Cfg        *Config
	Db         *DatabaseConfig
	Location   *time.Location
	FormatTime string
)

func init() {
	initEnv()
	initConfig(ENV.ReleaseMode)
	initTimeConfig()
}

func initEnv() {
	v := viper.New()
	v.SetConfigFile(envPath)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading Env file : %v", err))
	}

	ENV = &Env{
		AppName:     v.GetString("APP_NAME"),
		AppAuthor:   v.GetString("APP_AUTHOR"),
		AppVersion:  v.GetString("APP_VERSION"),
		AppHost:     v.GetString("APP_HOST"),
		AppPort:     v.GetString("APP_PORT"),
		ReleaseMode: v.GetString("RELEASE_MODE"),
	}
}

func initConfig(releaseMode string) {
	v := viper.New()
	v.SetConfigFile(configMap[releaseMode])
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config file : %v", err))
	}

	Cfg = &Config{
		Timezone:                   v.GetString("Timezone"),
		SecretKey:                  v.GetString("SecretKey"),
		TokenExpire:                v.GetInt("TokenExpire"),
		RefreshTokenExpire:         v.GetInt("RefreshTokenExpire"),
		FormatTime:                 v.GetString("FormatTime"),
		FormatDate:                 v.GetString("FormatDate"),
		MaxTimeoutGracefulShutdown: v.GetInt("MaxTimeoutGracefulShutdown"),
		Log: LogConfig{
			Path:       v.GetString("Log.Path"),
			Level:      v.GetString("Log.Level"),
			MaxSizeMB:  v.GetInt("Log.MaxSizeMB"),
			MaxBackups: v.GetInt("Log.MaxBackups"),
			MaxAgeDays: v.GetInt("Log.MaxAgeDays"),
		},
		CorsAllowOrigins: v.GetStringSlice("CorsAllowOrigins"),
	}

	Db = &DatabaseConfig{
		Type:         v.GetString("Database.Type"),
		Host:         v.GetString("Database.Host"),
		Port:         v.GetString("Database.Port"),
		User:         v.GetString("Database.User"),
		Password:     v.GetString("Database.Password"),
		Database:     v.GetString("Database.Database"),
		MaxOpenConns: v.GetInt("Database.MaxOpenConns"),
		MaxIdleConns: v.GetInt("Database.MaxIdleConns"),
		MaxLifetime:  v.GetInt("Database.MaxLifetime"),
		MaxIdleTime:  v.GetInt("Database.MaxIdleTime"),
	}

	Cfg.Database = *Db
}

func initTimeConfig() {
	loc, err := time.LoadLocation(Cfg.Timezone)
	if err != nil {
		panic(fmt.Sprintf("Failed to load location: %s", err.Error()))
	}
	Location = loc
	FormatTime = Cfg.FormatTime
}
