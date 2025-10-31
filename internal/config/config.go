package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host                 string
	Port                 string
	LimitCountPerRequest float64
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string
	DSN      string
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
	LogMode  bool
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// Load loads configuration from file
func Load(configPath ...string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.port", "8081")
	v.SetDefault("server.limitCountPerRequest", 10)
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.dbname", "ugin")
	v.SetDefault("database.logmode", true)
	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.accessTokenExpireDuration", 1)
	v.SetDefault("jwt.refreshTokenExpireDuration", 24)

	// Set config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if len(configPath) > 0 {
		v.AddConfigPath(configPath[0])
	} else {
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Allow environment variables to override config
	v.AutomaticEnv()

	cfg := &Config{}

	// Server config
	cfg.Server.Host = v.GetString("server.host")
	cfg.Server.Port = v.GetString("server.port")
	cfg.Server.LimitCountPerRequest = v.GetFloat64("server.limitCountPerRequest")

	// Database config
	cfg.Database.Driver = v.GetString("database.driver")
	cfg.Database.Database = v.GetString("database.dbname")
	cfg.Database.Username = v.GetString("database.username")
	cfg.Database.Password = v.GetString("database.password")
	cfg.Database.Host = v.GetString("database.host")
	cfg.Database.Port = v.GetString("database.port")
	cfg.Database.SSLMode = v.GetString("database.sslmode")
	cfg.Database.LogMode = v.GetBool("database.logmode")

	// Build DSN based on driver
	cfg.Database.DSN = buildDSN(cfg.Database)

	// JWT config
	cfg.JWT.Secret = v.GetString("server.secret")
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "change-me-in-production"
	}

	accessTokenHours := v.GetInt("server.accessTokenExpireDuration")
	if accessTokenHours == 0 {
		accessTokenHours = 1
	}
	cfg.JWT.AccessTokenDuration = time.Hour * time.Duration(accessTokenHours)

	refreshTokenHours := v.GetInt("server.refreshTokenExpireDuration")
	if refreshTokenHours == 0 {
		refreshTokenHours = 24
	}
	cfg.JWT.RefreshTokenDuration = time.Hour * time.Duration(refreshTokenHours)

	return cfg, nil
}

func buildDSN(db DatabaseConfig) string {
	switch db.Driver {
	case "sqlite":
		if db.Database == "" {
			return "ugin.db"
		}
		return db.Database + ".db"

	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.Username, db.Password, db.Host, db.Port, db.Database)

	case "postgres":
		sslmode := db.SSLMode
		if sslmode == "" {
			sslmode = "disable"
		}
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			db.Host, db.Port, db.Username, db.Password, db.Database, sslmode)

	default:
		return ""
	}
}

