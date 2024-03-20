package models

import (
	"errors"
	"math/rand"

	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

// User encapsulates user activities
type User interface {
	GetId() string
	GetApp() App
}

func NewUser(id string, app App) User {
	return &UserImpl{id: id, app: app}
}

type UserImpl struct {
	id  string
	app App
}

func (u *UserImpl) GetId() string {
	return u.id
}

func (u *UserImpl) GetApp() App {
	return u.app
}

func (u *UserImpl) chooseReactionIdToAdd(e Entity) (string, error) {
	var options []string

	// if max unique reactions is not reached for given entity
	// then any reaction can be added
	// otherwise only reactions that are already added to the given enetity (by someone)
	if len(getUniqReactionIds(e)) < e.GetNamespace().GetMaxUniqueReactions() {
		options = e.GetNamespace().GetAvailableReactionIds()
	} else {
		options = getUniqReactionIds(e)
	}

	// remove reactions that user already added
	utils.FilterIn(&options, func(r string) bool { return utils.Contains(e.GetMyReactionIds(), r) })

	if len(options) == 0 {
		return "", errors.New("has no options for reaction")
	}
	return options[rand.Intn(len(options))], nil
}

func getUniqReactionIds(e Entity) []string {
	var res []string
	for _, rc := range e.GetReactionsCount() {
		res = append(res, rc.ReactionId)
	}
	return res
}
