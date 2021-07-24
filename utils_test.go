package mongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlattenedMapFromInterface(t *testing.T) {
	type bar struct {
		Action string `json:"action"`
	}
	type dummy struct {
		Name string   `json:"name"`
		Bar  bar      `json:"bar"`
		Bars []string `json:"bars"`
	}

	tests := []struct {
		dummy       dummy
		expectedRes map[string]interface{}
	}{
		{
			dummy: dummy{
				Name: "dummy",
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
				"bars.0":     "tata",
				"bars.1":     "toto",
			},
		},
		{
			dummy: dummy{
				Name: "dummy",
				Bar: bar{
					Action: "dumb",
				},
				Bars: []string{"tata", "toto"},
			},
			expectedRes: map[string]interface{}{
				"name":       "dummy",
				"bar.action": "dumb",
				"bars.0":     "tata",
				"bars.1":     "toto",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, FlattenedMapFromInterface(test.dummy), test.expectedRes)
	}
}
