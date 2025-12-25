package config

type Configuration struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`

	Application struct {
		Handler         []string `yaml:"handler"`
		AllowedHost     string   `yaml:"allowed-host"`
		HTMLDirectory   string   `yaml:"html-dir"`
		WWWDirectory    string   `yaml:"www-dir"`
		AssetsDirectory string   `yaml:"assets-dir"`
	} `yaml:"application"`

	Logging struct {
		File  string `yaml:"file"`
		Level string `yaml:"level"`
	} `yaml:"logging"`
}
