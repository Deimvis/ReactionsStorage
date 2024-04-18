package rs

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/src/models"
	"github.com/Deimvis/reactionsstorage/tests/simulation/metrics"
)

type ClientType string

const (
	HTTP_CLIENT ClientType = "http"
	// GRPC_CLIENT = "grpc"
)

// thread-safe
type Client interface {
	GetReactions(req *models.ReactionsGETRequest) (models.Response, error)
	AddReaction(req *models.ReactionsPOSTRequest) (models.Response, error)
	RemoveReaction(req *models.ReactionsDELETERequest) (models.Response, error)

	SetConfiguration(req *models.ConfigurationPOSTRequest) (models.Response, error)
	GetNamespace(req *models.NamespaceGETRequest) (models.Response, error)
	GetAvailableReactions(req *models.AvailableReactionsGETRequest) (models.Response, error)
}

func NewClient(ct ClientType, host string, port int, ssl bool, logger *zap.SugaredLogger, recorder metrics.HTTPRecorder) (Client, error) {
	switch ct {
	case HTTP_CLIENT:
		return NewClientHTTP(host, port, ssl, logger, recorder), nil
	}
	return nil, fmt.Errorf("got unsupported client type: %s", ct)
}

func Expect[T models.Response](resp models.Response, err error) (T, error) {
	var respFinal T
	if err != nil {
		return respFinal, err
	}
	respFinal, ok := resp.(T)
	if !ok {
		return respFinal, fmt.Errorf("response has wrong type: %T", resp)
	}
	return respFinal, nil
}
