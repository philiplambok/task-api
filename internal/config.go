package internal

type Config struct {
	HTTPServer HTTPServerConfig `mapstructure:"http_server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
}

type HTTPServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	ExpirationHours   int    `mapstructure:"expiration_hours"`
}
