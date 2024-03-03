package services

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/storages"
)

func NewConfigurationService(lc fx.Lifecycle, cs *storages.ConfigurationStorage) *ConfigurationService {
	return &ConfigurationService{cs: cs}
}

type ConfigurationService struct {
	cs *storages.ConfigurationStorage
}

func (s *ConfigurationService) GetAvailableReactions(ctx context.Context, req *models.AvailableReactionsGETRequest) models.Response {
	if !s.cs.HasNamespace(ctx, req.Query.NamespaceId) {
		return &models.AvailableReactionsGETResponse404{Error: fmt.Sprintf("Namespace `%s` not found", req.Query.NamespaceId)}
	}
	reactions := s.cs.GetAvailableReactionsStrict(ctx, req.Query.NamespaceId)
	return &models.AvailableReactionsGETResponse200{Reactions: reactions}
}

// func (cs *ConfigurationService) SetConfiguration(ctx context.Context, conf *models.Configuration)
