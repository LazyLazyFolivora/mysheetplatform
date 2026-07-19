package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Alipay   AlipayConfig   `mapstructure:"alipay"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Email    EmailConfig    `mapstructure:"email"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	TTL    string `mapstructure:"ttl"`
}

type AlipayConfig struct {
	AppID          string `mapstructure:"app_id"`
	AppPrivateKey  string `mapstructure:"app_private_key"`
	AlipayPublicKey string `mapstructure:"alipay_public_key"`
	NotifyUrl      string `mapstructure:"notify_url"`
	ReturnUrl      string `mapstructure:"return_url"`
	SignType       string `mapstructure:"sign_type"`
	Charset        string `mapstructure:"charset"`
	GatewayUrl     string `mapstructure:"gateway_url"`
	Sandbox        bool   `mapstructure:"sandbox"`
}

type UploadConfig struct {
	Dir     string `mapstructure:"dir"`
	MaxSize int    `mapstructure:"max_size"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func Init(configPath string) *Config {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("SHEET")

	if err := v.ReadInConfig(); err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		zap.L().Fatal("failed to unmarshal config", zap.Error(err))
	}

	return &cfg
}
