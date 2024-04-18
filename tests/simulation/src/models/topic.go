package models

import (
	"fmt"

	"github.com/Deimvis/reactionsstorage/tests/simulation/src/utils"
)

type Topic interface {
	GetId() string
	GetEntities() []Entity
	CopyForUser() Topic
}

func NewTopic(id string, namespace Namespace, size int, shufflePerUser bool) Topic {
	var topic TopicImpl
	topic.id = id
	for i := 0; i < size; i++ {
		e := NewEntity(fmt.Sprintf("%s/%s", id, fmt.Sprint(i)), namespace)
		topic.entities = append(topic.entities, e)
	}
	return &topic
}

type TopicImpl struct {
	id             string
	entities       []Entity
	shufflePerUser bool
}

func (t *TopicImpl) GetId() string {
	return t.id
}

func (t *TopicImpl) GetEntities() []Entity {
	return t.entities
}

func (t *TopicImpl) CopyForUser() Topic {
	tc := t.copy()
	if t.shufflePerUser {
		utils.ShuffleIn(&tc.entities)
	}
	return tc
}

func (t *TopicImpl) copy() *TopicImpl {
	var topic TopicImpl
	topic.id = t.id
	for _, e := range t.entities {
		ecopy := NewEntity(e.GetId(), e.GetNamespace())
		topic.entities = append(topic.entities, ecopy)
	}
	return &topic
}
