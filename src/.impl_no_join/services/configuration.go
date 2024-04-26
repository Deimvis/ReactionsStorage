package services

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/src/storages"
	"github.com/Deimvis/reactionsstorage/src/utils"
)

func NewConfigurationService(lc fx.Lifecycle, cs *storages.ConfigurationStorage) *ConfigurationService {
	return &ConfigurationService{cs: cs}
}

type ConfigurationService struct {
	cs *storages.ConfigurationStorage
}

func (s *ConfigurationService) SetConfiguration(ctx context.Context, req *models.ConfigurationPOSTRequest) models.Response {
	err := utils.UntilFirstErr(
		func() error { return models.CheckCorrectness(&req.Body) },
		func() error { return s.cs.SetConfiguration(ctx, &req.Body) },
	)
	if err != nil {
		return &models.ConfigurationPOSTResponse422{Error: err.Error()}
	}
	return &models.ConfigurationPOSTResponse200{Status: "ok"}
}

func (s *ConfigurationService) GetNamespace(ctx context.Context, req *models.NamespaceGETRequest) models.Response {
	namespace, err := s.cs.GetNamespace(ctx, req.Query.NamespaceId)
	if namespace == nil {
		return &models.NamespaceGETResponse404{Error: fmt.Sprintf("Namespace `%s` not found: %s", req.Query.NamespaceId, err)}
	}
	return &models.NamespaceGETResponse200{Namespace: *namespace}
}

func (s *ConfigurationService) GetAvailableReactions(ctx context.Context, req *models.AvailableReactionsGETRequest) models.Response {
	ctx = storages.CtxAcquireConn(ctx, s.cs)
	defer storages.CtxReleaseConn(&ctx)

	if !s.cs.HasNamespace(ctx, req.Query.NamespaceId) {
		return &models.AvailableReactionsGETResponse404{Error: fmt.Sprintf("Namespace `%s` not found", req.Query.NamespaceId)}
	}
	reactions := s.cs.GetAvailableReactionsStrict(ctx, req.Query.NamespaceId)
	return &models.AvailableReactionsGETResponse200{Reactions: reactions}
}
