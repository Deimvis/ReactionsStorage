package config

import (
	"fmt"
	"os"

	"go.uber.org/fx"
	"gopkg.in/yaml.v3"

	"github.com/Deimvis/reactionsstorage/src/utils"
)

type Config struct {
	Gin Gin `yaml:"gin"`
}

func NewConfig(filePath *string) func(lc fx.Lifecycle) *Config {
	return func(lc fx.Lifecycle) *Config {
		cfg := &Config{}
		fileData := utils.Must(os.ReadFile(*filePath))
		fmt.Println(string(fileData))
		utils.Must0(yaml.Unmarshal(fileData, cfg))
		return cfg
	}
}

type Gin struct {
	General     GinGeneral     `yaml:"general"`
	Middlewares GinMiddlewares `yaml:"middlewares"`
	Handlers    GinHandlers    `yaml:"handlers"`
}

type GinGeneral struct {
	Mode           *string  `yaml:"mode"`
	TrustedProxies []string `yaml:"trusted_proxies"`
}

type GinMiddlewares struct {
	Logger     Option               `yaml:"logger"`
	Recovery   Option               `yaml:"recovery"`
	Prometheus PrometheusMiddleware `yaml:"prometheus"`
}

type PrometheusMiddleware struct {
	Option      `yaml:",inline"`
	MetricsPath string `yaml:"metrics_path"`
}

type GinHandlers struct {
	DebugHandlers GinDebugHandlers `yaml:"debug"`
}

type GinDebugHandlers struct {
	Pprof    PprofHandler    `yaml:"pprof"`
	MemUsage MemUsageHandler `yaml:"mem_usage"`
}

type PprofHandler struct {
	Option     `yaml:",inline"`
	PathPrefix *string `yaml:"path_prefix,omitempty"`
}

type MemUsageHandler struct {
	Option `yaml:",inline"`
	Path   *string `yaml:"path"`
}

type Option struct {
	Enabled bool `yaml:"enabled"`
}
