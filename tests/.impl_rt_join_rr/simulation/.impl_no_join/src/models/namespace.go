package models

import (
	"github.com/Deimvis/reactionsstorage/src/models"
	rs "github.com/Deimvis/reactionsstorage/tests/simulation/src/rs_client"
	"github.com/Deimvis/reactionsstorage/tests/simulation/src/utils"
)

type Namespace interface {
	GetId() string
	GetAvailableReactionIds() []string
	GetMaxUniqueReactions() int
	GetMutuallyExclusiveReactions() [][]string
}

func NewNamespace(namespaceId string, client rs.Client) Namespace {
	var namespace NamespaceImpl

	var req models.NamespaceGETRequest
	req.Query.NamespaceId = namespaceId
	resp, err := rs.Expect[*models.NamespaceGETResponse200](client.GetNamespace(&req))
	if err != nil {
		panic(err)
	}
	namespace.id = resp.Namespace.Id
	namespace.maxUniqReactions = resp.Namespace.MaxUniqReactions
	namespace.mutExclReactions = resp.Namespace.MutuallyExclusiveReactions

	var req2 models.AvailableReactionsGETRequest
	req2.Query.NamespaceId = namespaceId
	resp2, err := rs.Expect[*models.AvailableReactionsGETResponse200](client.GetAvailableReactions(&req2))
	if err != nil {
		panic(err)
	}
	namespace.availableReactionIds = utils.Map(resp2.Reactions, func(r models.Reaction) string { return r.Id })

	return &namespace
}

type NamespaceImpl struct {
	id                   string
	maxUniqReactions     int
	mutExclReactions     [][]string
	availableReactionIds []string
}

func (n *NamespaceImpl) GetId() string {
	return n.id
}

func (n *NamespaceImpl) GetAvailableReactionIds() []string {
	return n.availableReactionIds
}

func (n *NamespaceImpl) GetMaxUniqueReactions() int {
	return n.maxUniqReactions
}

func (n *NamespaceImpl) GetMutuallyExclusiveReactions() [][]string {
	return n.mutExclReactions
}

func GetConflictingReactionIds(n Namespace, reactionId string) []string {
	var res []string
	for _, conflictingGroup := range n.GetMutuallyExclusiveReactions() {
		if utils.Contains(conflictingGroup, reactionId) {
			otherIds := utils.Filter(conflictingGroup, func(id string) bool { return id != reactionId })
			for _, id := range otherIds {
				if !utils.Contains(res, id) {
					res = append(res, id)
				}
			}
		}
	}
	return res
}
