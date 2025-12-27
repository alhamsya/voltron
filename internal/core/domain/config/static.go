package config

import "time"

type Static struct {
	Env             string              `mapstructure:"env" validated:"required"`
	Frontend        *Frontend           `mapstructure:"frontend"`
	App             *App                `mapstructure:"app" validated:"required"`
	ServiceSpecific map[string]DBStatic `mapstructure:"service-specific" validate:"required"`
}

type Frontend struct {
	URL string `mapstructure:"url"`
}

type App struct {
	Rest Rest `mapstructure:"rest"`
}

type Rest struct {
	Port        int           `mapstructure:"port" validated:"required"`
	ReadTimeout time.Duration `mapstructure:"read-timeout"`
	IdleTimeout time.Duration `mapstructure:"idle-timeout"`
	Limiter     Limiter       `mapstructure:"limiter"`
}

type Limiter struct {
	Max        int           `mapstructure:"max"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type DBStatic struct {
	Main    DBConnStatic `mapstructure:"main" validate:"required"`
	Replica DBConnStatic `mapstructure:"replica" validate:"required"`
}

type DBConnStatic struct {
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required"`
	DBname          string `mapstructure:"dbname" validate:"required"`
	MaxOpenConns    int    `mapstructure:"max-open-conns"`
	MaxIdleConns    int    `mapstructure:"max-idle-conns"`
	ConnMaxLifetime int    `mapstructure:"conn-max-lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn-max-idle-time"`
}
