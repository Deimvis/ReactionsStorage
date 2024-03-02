package services

import (
	"go.uber.org/fx"

	"github.com/Deimvis/reactionsstorage/src/storages"
)

func NewConfigurationService(lc fx.Lifecycle, cs *storages.ConfigurationStorage) *ConfigurationService {
	return &ConfigurationService{cs: cs}
}

type ConfigurationService struct {
	cs *storages.ConfigurationStorage
}

// func (cs *ConfigurationService) SetConfiguration(ctx context.Context, conf *models.Configuration)
