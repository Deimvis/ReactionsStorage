package configs

import (
	"bytes"
	"os"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/Deimvis/reactionsstorage/src/utils"
)

type ServerConfig struct {
	Gin Gin `yaml:"gin"`
	PG  PG  `yaml:"pg"`
}

func NewServerConfig(filePath *string) func(lc fx.Lifecycle, logger *zap.SugaredLogger) *ServerConfig {
	return func(lc fx.Lifecycle, logger *zap.SugaredLogger) *ServerConfig {
		cfg := &ServerConfig{}
		fileData := utils.Must(os.ReadFile(*filePath))
		logger.Infof("Config:\n%s", string(fileData))
		decoder := yaml.NewDecoder(bytes.NewReader(fileData))
		decoder.KnownFields(true)
		utils.Must0(decoder.Decode(&cfg))
		logger.Debugf("Parsed config:\n%s", spew.Sdump(cfg))
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
	MetricsPath string            `yaml:"metrics_path"`
	Metrics     PrometheusMetrics `yaml:"metrics"`
}

type PrometheusMetrics struct {
	Gin   Option `yaml:"gin"`
	SQL   Option `yaml:"sql"`
	Debug Option `yaml:"debug"`
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
