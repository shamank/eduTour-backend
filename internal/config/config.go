package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

type (
	Config struct {
		HTTP       HTTPConfig
		SMTP       SMTPConfig
		Postgres   PostgresConfig
		AuthConfig AuthConfig
	}

	HTTPConfig struct {
		Host               string
		Port               string
		ReadTimeOut        time.Duration
		WriteTimeOut       time.Duration
		MaxHeaderMegabytes int
	}

	SMTPConfig struct {
		Host     string
		Port     int
		User     string
		Passowrd string
	}

	PostgresConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
		SignedKey       string
	}
)

func Init(configDir string) (*Config, error) {
	if err := parseConfig(configDir); err != nil {
		return nil, err
	}

	readTimeoutDuration, err := time.ParseDuration(viper.GetString("http.readTimeout"))
	if err != nil {
		return nil, err
	}

	writeTimeoutDuration, err := time.ParseDuration(viper.GetString("http.writeTimeout"))
	if err != nil {
		return nil, err
	}

	accessTokenTTL, err := time.ParseDuration(viper.GetString("auth.jwt.accessTTL"))
	if err != nil {
		return nil, err
	}
	refreshTokenTTL, err := time.ParseDuration(viper.GetString("auth.jwt.refreshTTL"))
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTP: HTTPConfig{
			Host:               viper.GetString("http.host"),
			Port:               viper.GetString("http.port"),
			ReadTimeOut:        readTimeoutDuration,
			WriteTimeOut:       writeTimeoutDuration,
			MaxHeaderMegabytes: viper.GetInt("http.maxHeaderBytes"),
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("smtp.host"),
			Port:     viper.GetInt("smtp.port"),
			User:     viper.GetString("smtp.user"),
			Passowrd: os.Getenv("SMTP_PASSWORD"),
		},
		Postgres: PostgresConfig{
			Host:     viper.GetString("pg.host"),
			Port:     viper.GetString("pg.port"),
			User:     viper.GetString("pg.user"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   viper.GetString("pg.dbname"),
			SSLMode:  viper.GetString("pg.sslmode"),
		},
		AuthConfig: AuthConfig{
			JWT: JWTConfig{
				AccessTokenTTL:  accessTokenTTL,
				RefreshTokenTTL: refreshTokenTTL,
				SignedKey:       os.Getenv("JWT_SIGNED_KEY"),
			},
			PasswordSalt:           os.Getenv("PASSWORD_SALT"),
			VerificationCodeLength: viper.GetInt("auth.verificationCodeLength"),
		},
	}, nil
}

func parseConfig(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := godotenv.Load(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}
