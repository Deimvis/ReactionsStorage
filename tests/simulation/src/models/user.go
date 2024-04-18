package models

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/mroth/weightedrand"
	"go.uber.org/zap"

	"github.com/Deimvis/reactionsstorage/tests/simulation/src/configs"
	"github.com/Deimvis/reactionsstorage/tests/simulation/src/utils"
)

// User encapsulates user activities
type User interface {
	GetId() string
	GetApp() App
	DoRandomAction()

	CanSwitchTopic() bool
	SwitchTopic()

	CanScroll() bool
	Scroll()

	CanAddReaction() bool
	AddReaction()

	CanRemoveReaction() bool
	RemoveReaction()
}

func NewUser(id string, app App, probs configs.ActionProbs, logger *zap.SugaredLogger) User {
	return &UserImpl{id: id, app: app, probs: probs, logger: logger}
}

type Action = func(u User)

type UserImpl struct {
	id    string
	app   App
	probs configs.ActionProbs

	logger *zap.SugaredLogger
}

func (u *UserImpl) GetId() string {
	return u.id
}

func (u *UserImpl) GetApp() App {
	return u.app
}

func (u *UserImpl) DoRandomAction() {
	u.chooseAction()(u)
}

func (u *UserImpl) chooseAction() Action {
	var choices []weightedrand.Choice
	if u.CanSwitchTopic() {
		choices = append(choices, weightedrand.NewChoice(User.SwitchTopic, u.probs.SwitchTopic))
	}
	if u.CanScroll() {
		choices = append(choices, weightedrand.NewChoice(User.Scroll, u.probs.Scroll))
	}
	if u.CanAddReaction() {
		choices = append(choices, weightedrand.NewChoice(User.AddReaction, u.probs.AddReaction))
	}
	if u.CanRemoveReaction() {
		choices = append(choices, weightedrand.NewChoice(User.RemoveReaction, u.probs.RemoveReaction))
	}
	if len(choices) == 0 {
		u.logger.Warnf("User %s has no options what to do next - skip this turn")
		return func(u User) {}
	}
	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		panic(err)
	}
	action, ok := chooser.Pick().(Action)
	if !ok {
		panic(fmt.Errorf("failed to convert chosen option to action"))
	}
	u.logger.Infof("User %s chose action: %s", u.id, getActionName(action))
	return action
}

func (u *UserImpl) CanSwitchTopic() bool {
	return true
}

func (u *UserImpl) SwitchTopic() {
	topicIds := u.app.GetAvailableTopicIds()
	id := topicIds[rand.Intn(len(topicIds))]
	u.app.SwitchTopic(id, u.GetId()).Wait()
}

func (u *UserImpl) CanScroll() bool {
	return u.app.CanScroll()
}

func (u *UserImpl) Scroll() {
	waitable, err := u.app.Scroll(u.GetId())
	if err != nil {
		panic(err)
	}
	waitable.Wait()
}

func (u *UserImpl) CanAddReaction() bool {
	for _, e := range u.app.GetVisibleEntities() {
		if len(getAddReactionOptions(e)) > 0 {
			return true
		}
	}
	return false
}

func (u *UserImpl) AddReaction() {
	entities := u.app.GetVisibleEntities()
	e, err := chooseRandomEntity(entities, func(e Entity) bool {
		options := getAddReactionOptions(e)
		return len(options) > 0
	})
	if err != nil {
		u.logger.Warnw("Failed to choose random entity", "err", err,
			"action", "add_reaction", "user", u.GetId())
		return
		// panic(err)
	}
	options := getAddReactionOptions(e)
	reactionId := options[rand.Intn(len(options))]
	u.app.AddReaction(e, u.GetId(), reactionId).Wait()
}

func (u *UserImpl) CanRemoveReaction() bool {
	for _, e := range u.app.GetVisibleEntities() {
		if len(getRemoveReactionOptions(e)) > 0 {
			return true
		}
	}
	return false
}

func (u *UserImpl) RemoveReaction() {
	entities := u.app.GetVisibleEntities()
	e, err := chooseRandomEntity(entities, func(e Entity) bool {
		options := getRemoveReactionOptions(e)
		return len(options) > 0
	})
	if err != nil {
		u.logger.Warnw("Failed to choose random entity", "err", err,
			"action", "remove_reaction", "user", u.GetId())
		return
		// panic(err)
	}
	options := getRemoveReactionOptions(e)
	reactionId := options[rand.Intn(len(options))]
	u.app.RemoveReaction(e, u.GetId(), reactionId).Wait()
}

func chooseRandomEntity(entities []Entity, pred func(e Entity) bool) (Entity, error) {
	for _, e := range utils.Shuffle(entities) {
		if pred(e) {
			// Shuffle returns a copy of original slice
			// But underlying elements should be pointers
			utils.AssertPtr(e)
			return e, nil
		}
	}
	return nil, errors.New("no entity satisfied given predicate")
}

func getAddReactionOptions(e Entity) []string {
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
	options = utils.Filter(options, func(r string) bool { return utils.Contains(e.GetMyReactionIds(), r) })

	return options
}

func getRemoveReactionOptions(e Entity) []string {
	return e.GetMyReactionIds()
}

func getUniqReactionIds(e Entity) []string {
	var res []string
	for _, rc := range e.GetReactionsCount() {
		res = append(res, rc.ReactionId)
	}
	return res
}

func getActionName(a Action) string {
	fnName := utils.GetFnName(a)
	parts := strings.Split(fnName, ".")
	if len(parts) == 0 {
		panic(fmt.Errorf("action has bad function name: %s", fnName))
	}
	return parts[len(parts)-1]
}
