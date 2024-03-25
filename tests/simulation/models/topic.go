package models

import (
	"fmt"

	"github.com/Deimvis/reactionsstorage/tests/simulation/utils"
)

type Topic interface {
	GetId() string
	GetEntities() []Entity
	ShuffleEntities()
	Copy() Topic
}

func NewTopic(id string, size int, namespace Namespace) Topic {
	var topic TopicImpl
	topic.id = id
	for i := 0; i < size; i++ {
		e := NewEntity(fmt.Sprintf("%s/%s", id, fmt.Sprint(i)), namespace)
		topic.entities = append(topic.entities, e)
	}
	return &topic
}

type TopicImpl struct {
	id       string
	entities []Entity
}

func (t *TopicImpl) GetId() string {
	return t.id
}

func (t *TopicImpl) GetEntities() []Entity {
	return t.entities
}

func (t *TopicImpl) ShuffleEntities() {
	utils.ShuffleIn(&t.entities)
}

func (t *TopicImpl) Copy() Topic {
	var topic TopicImpl
	topic.id = t.id
	for _, e := range t.entities {
		ecopy := NewEntity(e.GetId(), e.GetNamespace())
		topic.entities = append(topic.entities, ecopy)
	}
	return &topic
}