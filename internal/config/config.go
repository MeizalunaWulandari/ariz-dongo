package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string
	AppEnv     string
	AppPort    string
	AppVersion string

	AllowSite string

	Fieldsa FieldsaConfig

	// Disiapkan untuk future use.
	DBDriver   string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

type FieldsaConfig struct {
	BaseURL         string
	DispatchBaseURL string
	Username        string
	Password        string
	LoginBearer     string
}

func Load() Config {
	// Membaca .env jika tersedia.
	// Tidak error jika file .env tidak ada.
	_ = godotenv.Load()

	return Config{
		AppName:    getEnv("APP_NAME", "ariz-dongo"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "3000"),
		AppVersion: getEnv("APP_VERSION", "1.0.0"),

		AllowSite: getEnv(
			"ALLOW_SITE",
			"",
		),

		Fieldsa: FieldsaConfig{
			BaseURL: getEnv(
				"FIELDSA_BASE_URL",
				"https://fsservice.fiberstar.co.id",
			),

			DispatchBaseURL: getEnv(
				"FIELDSA_DISPATCH_BASE_URL",
				"https://fieldsaapi.fiberstar.co.id",
			),

			Username: getEnv(
				"FIELDSA_USERNAME",
				"",
			),

			Password: getEnv(
				"FIELDSA_PASSWORD",
				"",
			),

			LoginBearer: getEnv(
				"FIELDSA_LOGIN_BEARER",
				"",
			),
		},

		DBDriver:   getEnv("DB_DRIVER", ""),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
