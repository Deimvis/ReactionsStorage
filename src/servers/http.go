package servers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/metrics"
	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/services"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func NewHTTPServer(lc fx.Lifecycle, cs *services.ConfigurationService, rs *services.ReactionsService, logger *zap.SugaredLogger) *http.Server {
	addr := fmt.Sprintf(":%s", utils.Getenv("PORT", "8080"))
	s := &http.Server{
		Addr:    addr,
		Handler: NewRouter(cs, rs, logger),
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

func NewRouter(cs *services.ConfigurationService, rs *services.ReactionsService, logger *zap.SugaredLogger) *gin.Engine {
	if utils.IsDebugEnv() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)
	UsePrometheusMiddleware(router)

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

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

	router.POST("/reactions/events", func(c *gin.Context) {
		// TODO: process consequently add, remove reaction events
		c.String(500, "TODO: implement")
	})

	return router
}

func UsePrometheusMiddleware(router *gin.Engine) {
	customMetrics := []*ginprometheus.Metric{
		metrics.GINReqDurV2Wrap.Metric,
		metrics.SQLReqCntWrap.Metric,
		metrics.SQLReqDurWrap.Metric,
	}
	p := ginprometheus.NewPrometheus("gin", customMetrics)
	p.Use(router)
	metrics.GINReqDurV2 = metrics.GINReqDurV2Wrap.Unwrap()
	metrics.SQLReqCnt = metrics.SQLReqCntWrap.Unwrap()
	metrics.SQLReqDur = metrics.SQLReqDurWrap.Unwrap()

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
