package config

type Config struct {
	defaultNamespace string `yaml:"default_namespace"`
}

func NewConfig() Config {
	return Config{defaultNamespace: "default"}
}
