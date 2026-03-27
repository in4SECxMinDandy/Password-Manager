package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Crypto    CryptoConfig
}

type ServerConfig struct {
	Address       string
	AllowOrigins  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	AccessExpiry    int // minutes
	RefreshExpiry   int // days
}

type CryptoConfig struct {
	Argon2Memory      uint32
	Argon2Iterations uint32
	Argon2Parallelism uint8
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("server.allowOrigins", "*")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "passwordmanager")
	viper.SetDefault("database.password", "passwordmanager_secret")
	viper.SetDefault("database.dbName", "passwordmanager")
	viper.SetDefault("database.sslMode", "disable")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("jwt.accessExpiry", 15)
	viper.SetDefault("jwt.refreshExpiry", 7)
	viper.SetDefault("crypto.argon2Memory", 65536)
	viper.SetDefault("crypto.argon2Iterations", 3)
	viper.SetDefault("crypto.argon2Parallelism", 4)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
