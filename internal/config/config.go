package config

import (
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ServerCfg *ServerCfg
	WebCfg    *WebCfg
	DbCfg     *DbCfg
}

type ServerCfg struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type WebCfg struct {
	TemplatesDir   string
	StaticFilesDir string
}

type DbCfg struct {
	DbPath          string
	DbConfigs       string
	DbMigrationPath string
}

func GetConfig() *Config {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	baseDir := filepath.Dir(exePath)

	return &Config{
		ServerCfg: &ServerCfg{
			Addr:         ":8080",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		WebCfg: &WebCfg{
			TemplatesDir:   "internal/web-storage/templates",
			StaticFilesDir: "internal/web-storage/static",
		},
		DbCfg: &DbCfg{
			DbPath:          filepath.Join(baseDir, "internal", "db", "app.db"),
			DbConfigs:       "?_foreign_keys=on",
			DbMigrationPath: "internal/db/migrations/up.sql",
		},
	}
}
