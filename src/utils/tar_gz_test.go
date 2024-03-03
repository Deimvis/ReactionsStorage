package utils_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Deimvis/reactionsstorage/src/utils"
)

type Obj1 struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Obj2 struct {
	Field1 int    `yaml:"field_1"`
	Field2 string `yaml:"field_2"`
}

type YamlGroup struct {
	Yaml1 []Obj1 `filename:"1.yaml"`
	Yaml2 Obj2   `filename:"2.yaml"`
}

func TestCorrectness(t *testing.T) {
	yg := YamlGroup{
		Yaml1: []Obj1{
			{
				Key:   "key",
				Value: "value",
			},
			{
				Key:   "a",
				Value: "b",
			},
		},
		Yaml2: Obj2{
			Field1: 42,
			Field2: "answer",
		},
	}
	tar, err := utils.CreateTarGz(yg)
	if err != nil {
		t.Fatalf("failed to create tar gz: %w", err)
	}
	var ygCopy YamlGroup
	err = utils.ExtractTarGz(tar, &ygCopy)
	if err != nil {
		t.Fatalf("failed to extract tar gz: %w", err)
	}
	require.True(t, reflect.DeepEqual(yg, ygCopy))
}
