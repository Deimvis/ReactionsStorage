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

func NewHTTPServer(lc fx.Lifecycle, reactionsService *services.ReactionsService) *http.Server {
	s := &http.Server{
		Addr:    ":8080",
		Handler: NewRouter(reactionsService),
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

func NewRouter(reactionsService *services.ReactionsService) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Ok"})
	})

	router.GET("/reactions", func(c *gin.Context) {
		var req models.ReactionsGETRequest
		req.Query.NamespaceId = c.Query("namespace_id")
		req.Query.EntityId = c.Query("entity_id")
		req.Query.UserId = c.Query("user_id")
		log.Println("Process request:", req)

		resp := reactionsService.GetUserReactions(c, req)
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

		resp := reactionsService.AddUserReaction(c, req)
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

		resp := reactionsService.RemoveUserReaction(c, req)
		c.JSON(resp.Code(), resp)

		log.Println("Return response:", resp)
	})

	router.POST("/reactions/events", func(c *gin.Context) {
		fmt.Println("POST /reactions/events was called")
		// TODO: process consequently add,remove reaction events
	})

	return router
}
