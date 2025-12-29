package modelConfig

import "time"

type Static struct {
	Env             string              `mapstructure:"env" validated:"required"`
	Frontend        *Frontend           `mapstructure:"frontend"`
	App             *App                `mapstructure:"app" validated:"required"`
	ServiceSpecific map[string]DBStatic `mapstructure:"service-specific" validate:"required"`
	RabbitMQ        RabbitMQStatic      `mapstructure:"rabbit-mq"`
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
	Primary DBConnStatic `mapstructure:"primary" validate:"required"`
	Replica DBConnStatic `mapstructure:"replica" validate:"required"`
}

type DBConnStatic struct {
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required"`
	Name            string `mapstructure:"name" validate:"required"`
	MaxConns        int    `mapstructure:"max-conns"`
	MinIdleConns    int    `mapstructure:"min-idle-conns"`
	MaxConnLifetime int    `mapstructure:"max-conn-lifetime"`
	MaxConnIdleTime int    `mapstructure:"max-conn-idle-time"`
}

type RabbitMQStatic struct {
	Host string `mapstructure:"host" validate:"required"`
	Port int    `mapstructure:"port" validate:"required"`
}
