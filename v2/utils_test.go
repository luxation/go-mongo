package v2

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestFlattenedMapFromInterface(t *testing.T) {
	type bar struct {
		Action string `json:"action"`
	}
	type dummy struct {
		Name *string  `json:"name"`
		Bar  bar      `json:"bar"`
		Bars []string `json:"bars"`
	}

	dummyName := "dummy"

	tests := []struct {
		dummy       dummy
		expectedRes map[string]interface{}
	}{
		{
			dummy: dummy{
				Name: &dummyName,
				Bar: bar{
					Action: "dumb",
				},
				Bars: nil,
			},
			expectedRes: map[string]interface{}{
				"name":       "dummy",
				"bar.action": "dumb",
			},
		},
		{
			dummy: dummy{
				Bar: bar{
					Action: "dumb",
				},
				Bars: []string{"tata", "toto"},
			},
			expectedRes: map[string]interface{}{
				"bar.action": "dumb",
				"bars":       []interface{}{"tata", "toto"},
			},
		},
		{
			dummy: dummy{
				Name: &dummyName,
				Bar: bar{
					Action: "dumb",
				},
				Bars: []string{"tata", "toto"},
			},
			expectedRes: map[string]interface{}{
				"name":       "dummy",
				"bar.action": "dumb",
				"bars":       []interface{}{"tata", "toto"},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedRes, FlattenedMapFromInterface(test.dummy))
	}
}

func TestFlattenedMapFromBson(t *testing.T) {
	obj := bson.M{"test": "foo", "test2": nil}

	assert.Equal(t, map[string]interface{}{"test": "foo", "test2": nil}, FlattenedMapFromInterface(obj))
}
