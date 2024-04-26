package servers

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/configs"
	"github.com/Deimvis/reactionsstorage/src/metrics"
	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func NewHTTPServer(lc fx.Lifecycle, cfg *configs.ServerConfig, cs *services.ConfigurationService, rs *services.ReactionsService, logger *zap.SugaredLogger) *http.Server {
	addr := fmt.Sprintf(":%s", utils.Getenv("PORT", "8080"))
	s := &http.Server{
		Addr:    addr,
		Handler: NewRouter(&cfg.Gin, cs, rs, logger),
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatalf("listen: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
	return s
}

func NewRouter(cfg *configs.Gin, cs *services.ConfigurationService, rs *services.ReactionsService, logger *zap.SugaredLogger) *gin.Engine {
	router := InitRouter(&cfg.General)
	UseMiddlewares(&cfg.Middlewares, router)
	SetHandlers(&cfg.Handlers, router, cs, rs, logger)
	return router
}

func InitRouter(cfg *configs.GinGeneral) *gin.Engine {
	if utils.IsDebugEnv() {
		gin.SetMode(gin.DebugMode)
	} else if cfg.Mode != nil {
		gin.SetMode(*cfg.Mode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.SetTrustedProxies(cfg.TrustedProxies)
	return router
}

func UseMiddlewares(cfg *configs.GinMiddlewares, router *gin.Engine) {
	if cfg.Logger.Enabled {
		router.Use(gin.Logger())
	}
	if cfg.Recovery.Enabled {
		router.Use(gin.Recovery())
	}
	if cfg.Prometheus.Enabled {
		UsePrometheusMiddleware(cfg, router)
	}
}

func UsePrometheusMiddleware(cfg *configs.GinMiddlewares, router *gin.Engine) {
	customMetrics := []*ginprometheus.Metric{
		metrics.GINReqDurV2Wrap.Metric,
		metrics.SQLReqCntWrap.Metric,
		metrics.SQLReqDurWrap.Metric,
		metrics.GetEntityReactionsCountWrap.Metric,
		metrics.GetUniqEntityUserReactionsWrap.Metric,
		metrics.GetEntityReactionsCountQueryWrap.Metric,
		metrics.GetEntityReactionsCountCollectRowsWrap.Metric,
		metrics.GETReactionsAcquireWrap.Metric,
	}

	p := ginprometheus.NewPrometheus("gin", customMetrics)
	p.MetricsPath = cfg.Prometheus.MetricsPath
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		return c.Request.URL.Path
	}
	p.Use(router)

	metrics.GINReqDurV2 = metrics.GINReqDurV2Wrap.Unwrap()
	metrics.SQLReqCnt = metrics.SQLReqCntWrap.Unwrap()
	metrics.SQLReqDur = metrics.SQLReqDurWrap.Unwrap()
	metrics.GetEntityReactionsCount = metrics.GetEntityReactionsCountWrap.Unwrap()
	metrics.GetUniqEntityUserReactions = metrics.GetUniqEntityUserReactionsWrap.Unwrap()
	metrics.GetEntityReactionsCountQuery = metrics.GetEntityReactionsCountQueryWrap.Unwrap()
	metrics.GetEntityReactionsCountCollectRows = metrics.GetEntityReactionsCountCollectRowsWrap.Unwrap()
	metrics.GETReactionsAcquire = metrics.GETReactionsAcquireWrap.Unwrap()

	// Record request duration with GINReqDurV2
	router.Use(func(c *gin.Context) {
		if c.Request.URL.String() == p.MetricsPath {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		elapsed := float64(time.Since(start)) / float64(time.Second)
		status := strconv.Itoa(c.Writer.Status())
		url := p.ReqCntURLLabelMappingFn(c)
		metrics.GINReqDurV2Wrap.Unwrap().WithLabelValues(status, c.Request.Method, url).Observe(elapsed)
	})
}

func SetHandlers(cfg *configs.GinHandlers, router *gin.Engine, cs *services.ConfigurationService, rs *services.ReactionsService, logger *zap.SugaredLogger) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/configuration", func(c *gin.Context) {
		var req models.ConfigurationPOSTRequest
		req.Headers = c.Request.Header

		switch c.ContentType() {
		case "application/yaml":
			err := c.BindYAML(&req.Body)
			if err != nil {
				logger.Infof("Bad request: %s\n", err)
				return // 400
			}
		case "application/json":
			err := c.BindJSON(&req.Body)
			if err != nil {
				logger.Infof("Bad request: %s\n", err)
				return // 400
			}
		case "application/gzip":
			c.String(500, "TODO: implement")
			return
		default:
			msg := fmt.Sprintf("unsupported Content-Type: %s", c.ContentType())
			resp := &models.ConfigurationPOSTResponse415{Error: msg}
			c.JSON(resp.Code(), resp)
			return
		}
		logger.Debugln("Process request:", req)

		resp := cs.SetConfiguration(c, &req)
		c.JSON(resp.Code(), resp)
	})

	router.GET("/configuration/namespace", func(c *gin.Context) {
		var req models.NamespaceGETRequest
		req.Query.NamespaceId = c.Query("namespace_id")
		logger.Debugln("Process request:", req)

		resp := cs.GetNamespace(c, &req)
		c.JSON(resp.Code(), resp)
	})

	router.GET("/configuration/available_reactions", func(c *gin.Context) {
		var req models.AvailableReactionsGETRequest
		req.Query.NamespaceId = c.Query("namespace_id")
		logger.Debugln("Process request:", req)

		resp := cs.GetAvailableReactions(c, &req)
		c.JSON(resp.Code(), resp)
	})

	router.GET("/reactions", func(c *gin.Context) {
		var req models.ReactionsGETRequest
		req.Query.NamespaceId = c.Query("namespace_id")
		req.Query.EntityId = c.Query("entity_id")
		req.Query.UserId = c.Query("user_id")
		logger.Debugln("Process request:", req)

		resp := rs.GetUserReactions(c, req)
		c.JSON(resp.Code(), resp)
	})

	router.POST("/reactions", func(c *gin.Context) {
		var req models.ReactionsPOSTRequest
		force := c.Query("force") == "true"
		req.Query.Force = &force
		err := c.BindJSON(&req.Body)
		if err != nil {
			logger.Infof("Bad request: %s\n", err)
			return // 400
		}
		logger.Debugln("Process request:", req)

		resp := rs.AddUserReaction(c, req)
		c.JSON(resp.Code(), resp)

		logger.Debugln(spew.Sprintf("Return response: %v", resp))
	})

	router.DELETE("/reactions", func(c *gin.Context) {
		var req models.ReactionsDELETERequest
		err := c.BindJSON(&req.Body)
		if err != nil {
			logger.Infof("Bad request:\n%s", err)
			return // 400
		}
		logger.Debugln("Process request:", req)

		resp := rs.RemoveUserReaction(c, req)
		c.JSON(resp.Code(), resp)

		logger.Debugln(spew.Sprintf("Return response: %v", resp))
	})

	SetDebugHandlers(&cfg.DebugHandlers, router)
}

func SetDebugHandlers(cfg *configs.GinDebugHandlers, router *gin.Engine) {
	if cfg.Pprof.Enabled {
		if cfg.Pprof.PathPrefix != nil {
			pprof.Register(router, *cfg.Pprof.PathPrefix)
		} else {
			pprof.Register(router)
		}
	}
	if cfg.MemUsage.Enabled {
		path := defaultMemUsagePath
		if cfg.MemUsage.Path != nil {
			path = *cfg.MemUsage.Path
		}
		router.GET(path, func(c *gin.Context) {
			MB := func(B uint64) uint64 { return B / 1024 / 1024 }
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			stat := struct {
				AllocMB      uint64 `json:"alloc_MB"`
				TotalAllocMB uint64 `json:"total_alloc_MB"`
				SysMB        uint64 `json:"sys_MB"`
				NumGC        uint32 `json:"num_gc"`
			}{
				AllocMB:      MB(m.Alloc),
				TotalAllocMB: MB(m.TotalAlloc),
				SysMB:        MB(m.Sys),
				NumGC:        m.NumGC,
			}
			c.JSON(200, stat)
		})
	}
}

var defaultMemUsagePath = "/debug/sys/mem"
