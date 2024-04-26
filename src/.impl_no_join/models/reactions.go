package models

type UserReaction struct {
	NamespaceId string `json:"namespace_id" binding:"required"`
	EntityId    string `json:"entity_id" binding:"required"`
	ReactionId  string `json:"reaction_id" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
}

type ReactionCount struct {
	ReactionId string `json:"reaction_id"`
	Count      int    `json:"count"`
}

type UserReactionsWithinEntity struct {
	UserId    string   `json:"user_id"`
	Reactions []string `json:"reactions"`
}

func (rc ReactionCount) FromMap(m map[string]int) []ReactionCount {
	res := []ReactionCount{}
	for k, v := range m {
		res = append(res, ReactionCount{ReactionId: k, Count: v})
	}
	return res
}

func (rc ReactionCount) ToMap(rcs []ReactionCount) map[string]int {
	res := make(map[string]int)
	for _, rc := range rcs {
		res[rc.ReactionId] = rc.Count
	}
	return res
}
