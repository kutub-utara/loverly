package config

import (
	"context"
	"fmt"
	"loverly/lib/log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"
)

/*
	All config should be required.
	Optional only allowed if zero value of the type is expected being the default value.
	time.Duration units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”. as in time.ParseDuration().
*/

type (
	Postgres struct {
		ConnURI            string        `mapstructure:"PG_CONN_URI" validate:"required"`
		MaxPoolSize        int           `mapstructure:"PG_MAX_POOL_SZE"` //Optional, default to 0 (zero value of int)
		MaxIdleConnections int           `mapstructure:"PG_MAX_IDLE_CONNECTIONS"`
		MaxIdleTime        time.Duration `mapstructure:"PG_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
		MaxLifeTime        time.Duration `mapstructure:"PG_MAX_LIFE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
	}

	PostgresReader struct {
		ConnURI            string        `mapstructure:"PG_READER_CONN_URI" validate:"required"`
		MaxPoolSize        int           `mapstructure:"PG_READER_MAX_POOL_SZE"` //Optional, default to 0 (zero value of int)
		MaxIdleConnections int           `mapstructure:"PG_READER_MAX_IDLE_CONNECTIONS"`
		MaxIdleTime        time.Duration `mapstructure:"PG_READER_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
		MaxLifeTime        time.Duration `mapstructure:"PG_READER_MAX_LIFE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
	}

	Redis struct {
		Host     string `mapstructure:"REDIS_HOST" validate:"required"`
		Password string `mapstructure:"REDIS_PASSWORD"`
	}

	Translation struct {
		FilePath            string   `mapstructure:"TRANSLATION_FILE_PATH"`
		LanguagePreferences []string `mapstructure:"TRANSLATION_LANG_PREFERENCES"`
		DefaultLanguage     string   `mapstructure:"TRANSLATION_DEAULT_LANG"`
	}

	JWTKey struct {
		KeyId     string `mapstructure:"JWK_KID" validate:"required"`
		SignKey   string `mapstructure:"ACCESS_TOKEN_RSA256_PRIVATE_KEY" validate:"required"` //RSA Private Key in PEM
		VerifyKey string `mapstructure:"ACCESS_TOKEN_RSA256_PUBLIC_KEY" validate:"required"`  //RSA Public Key in PEM
	}

	Configuration struct {
		ServiceName          string         `mapstructure:"SERVICE_NAME"`
		TraceEndpoint        string         `mapstructure:"TRACE_ENDPOINT"`
		TraceRate            float64        `mapstructure:"TRACE_RATE"`
		Postgres             Postgres       `mapstructure:",squash"`
		PostgresReader       PostgresReader `mapstructure:",squash"`
		Translation          Translation    `mapstructure:",squash"`
		Redis                Redis          `mapstructure:",squash"`
		TokenIssuer          string         `mapstructure:"TOKEN_ISSUER" validate:"required"`
		IatLeeway            time.Duration  `mapstructure:"IAT_LEEWAY" validate:"required"` //Leeway time for iat to accommodate server time discrepancy
		JWTKey               JWTKey         `mapstructure:",squash"`
		AccessTokenValidity  time.Duration  `mapstructure:"ACCESS_TOKEN_VALID_FOR" validate:"required"`
		RefreshTokenValidity time.Duration  `mapstructure:"REFRESH_TOKEN_VALID_FOR" validate:"required"`
		AuthorizationCode    time.Duration  `mapstructure:"AUTHZ_CODE_VALID_FOR" validate:"required"`

		Environment string `mapstructure:"ENV" validate:"required,oneof=development staging production"`
		BindAddress int    `mapstructure:"BIND_ADDRESS" validate:"required"`
		LogLevel    int    `mapstructure:"LOG_LEVEL" validate:"required"`
	}
)

func InitConfig(ctx context.Context, log log.Interface) (*Configuration, error) {
	var cfg Configuration

	viper.SetConfigType("env")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if !os.IsNotExist(err) {
		viper.SetConfigFile(envFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Error(ctx, fmt.Sprintf("failed to read config:%v", err))
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Error(ctx, fmt.Sprintf("failed to bind config:%v", err))
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Error(ctx, fmt.Sprintf("invalid config:%v", err))
		}
		log.Error(ctx, fmt.Sprintf("failed to load config"))
		return nil, err
	}

	log.Info(ctx, fmt.Sprintf("Config loaded: %+v", cfg))
	return &cfg, nil
}

func (cfg Translation) TranslationJSONFiles() []string {
	var files []string
	languages := append(cfg.LanguagePreferences, cfg.DefaultLanguage)
	for _, lang := range languages {
		fileName := fmt.Sprintf("%s.all.json", lang)
		files = append(files, filepath.Join(cfg.FilePath, fileName))
	}
	return files
}
