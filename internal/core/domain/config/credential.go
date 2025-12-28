package config

type Credential struct {
	ServiceSpecific map[string]DBCredential `mapstructure:"service-specific" validate:"required"`
}

type DBConnCredential struct {
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}

type DBCredential struct {
	Primary DBConnCredential `mapstructure:"primary" validate:"required"`
	Replica DBConnCredential `mapstructure:"replica" validate:"required"`
}
