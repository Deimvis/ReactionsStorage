package configs

import (
	"fmt"
	"os"

	"go.uber.org/fx"
	"gopkg.in/yaml.v3"

	"github.com/Deimvis/reactionsstorage/src/utils"
)

type ServerConfig struct {
	Gin Gin `yaml:"gin"`
	PG  PG  `yaml:"pg"`
}

func NewServerConfig(filePath *string) func(lc fx.Lifecycle) *ServerConfig {
	return func(lc fx.Lifecycle) *ServerConfig {
		cfg := &ServerConfig{}
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

type PG struct {
	Pool PGPool `yaml:"pool"`
}

type PGPool struct {
	MinConns               *int32 `yaml:"min_conns"`
	MaxConns               *int32 `yaml:"max_conns"`
	MaxConnLifetimeJitterS *int   `yaml:"max_conn_lifetime_jitter_s"`
}

type Option struct {
	Enabled bool `yaml:"enabled"`
}
