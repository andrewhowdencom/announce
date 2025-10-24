package config

// Email holds the configuration for the email client.
type Email struct {
	From     string `mapstructure:"from"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Config holds the application configuration.
type Config struct {
	Email Email `mapstructure:"email"`
}
