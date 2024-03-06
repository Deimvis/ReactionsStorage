package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/services"
)

func NewHTTPServer(lc fx.Lifecycle, cs *services.ConfigurationService, rs *services.ReactionsService) *http.Server {
	s := &http.Server{
		Addr:    ":8080",
		Handler: NewRouter(cs, rs),
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
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

func NewRouter(cs *services.ConfigurationService, rs *services.ReactionsService) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/configuration", func(c *gin.Context) {
		c.String(500, "TODO: implement")

		var resp models.Response
		switch c.ContentType() {
		case "application/yaml":

		default:
			msg := fmt.Sprintf("unsupported Content-Type: %s", c.ContentType())
			resp = &models.ConfigurationPOSTResponse415{Error: msg}
		}
		c.JSON(resp.Code(), resp)

	})
	router.GET("/configuration/namespace", func(c *gin.Context) {
		var req models.NamespaceGETRequest
		req.Query.NamespaceId = c.Query("namespace_id")
		log.Println("Process request:", req)

		resp := cs.GetNamespace(c, &req)
		c.JSON(resp.Code(), resp)
	})
	router.GET("/configuration/available_reactions", func(c *gin.Context) {
		var req models.AvailableReactionsGETRequest
		req.Query.NamespaceId = c.Query("namespace_id")
		log.Println("Process request:", req)

		resp := cs.GetAvailableReactions(c, &req)
		c.JSON(resp.Code(), resp)
	})

	router.GET("/reactions", func(c *gin.Context) {
		var req models.ReactionsGETRequest
		req.Query.NamespaceId = c.Query("namespace_id")
		req.Query.EntityId = c.Query("entity_id")
		req.Query.UserId = c.Query("user_id")
		log.Println("Process request:", req)

		resp := rs.GetUserReactions(c, req)
		c.JSON(resp.Code(), resp)
	})

	router.POST("/reactions", func(c *gin.Context) {
		var req models.ReactionsPOSTRequest
		force := c.Query("force") == "true"
		req.Query.Force = &force
		err := c.BindJSON(&req.Body)
		if err != nil {
			log.Printf("Bad request:\n%s", err)
			return // 400
		}
		log.Println("Process request:", req)

		resp := rs.AddUserReaction(c, req)
		c.JSON(resp.Code(), resp)

		log.Println("Return response:", resp)
	})

	router.DELETE("/reactions", func(c *gin.Context) {
		var req models.ReactionsDELETERequest
		err := c.BindJSON(&req.Body)
		if err != nil {
			log.Printf("Bad request:\n%s", err)
			return // 400
		}
		log.Println("Process request:", req)

		resp := rs.RemoveUserReaction(c, req)
		c.JSON(resp.Code(), resp)

		log.Println("Return response:", resp)
	})

	router.POST("/reactions/events", func(c *gin.Context) {
		fmt.Println("POST /reactions/events was called")
		// TODO: process consequently add,remove reaction events
	})

	return router
}
