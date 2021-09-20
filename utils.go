package mongo

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FlattenedMapFromInterface(from interface{}) map[string]interface{} {
	jsonFields := make(map[string]interface{})
	switch from.(type) {
	case primitive.M, primitive.D:
		marshaled, err := json.Marshal(from)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(marshaled, &jsonFields)
		if err != nil {
			panic(err)
		}

		return jsonFields
	default:
		resultFields := make(map[string]interface{})

		marshaled, err := json.Marshal(from)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(marshaled, &jsonFields)
		if err != nil {
			panic(err)
		}

		flattenNestMap("", jsonFields, resultFields)

		return resultFields
	}
}

func flattenNestMap(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flattenNestMap(prefix+k, child, dest)
		case nil:
			break
		default:
			dest[prefix+k] = v
		}
	}
}
