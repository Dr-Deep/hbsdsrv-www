package config

type Configuration struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`

	Application struct {
		Handler        []string `yaml:"handler"`
		AllowedOrigins []string `yaml:"allowed-origins"`
		HTMLDirectory  string   `yaml:"html-dir"`
		WWWDirectory   string   `yaml:"www-dir"`
	} `yaml:"application"`

	Logging struct {
		File  string `yaml:"file"`
		Level string `yaml:"level"`
	} `yaml:"logging"`
}
