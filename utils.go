package mongo

import (
	"encoding/json"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func FlattenedMapFromInterface(from interface{}) map[string]interface{} {
	jsonFields := make(map[string]interface{})
	switch from.(type) {
	case primitive.M, primitive.D, primitive.E, primitive.A, primitive.Regex:
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

		keys := make(map[string]interface{})

		jsonMarshaled, err := json.Marshal(from)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(jsonMarshaled, &keys)
		if err != nil {
			panic(err)
		}

		kk := make(map[string]string)

		fillCamelCaseFields(keys, kk)

		marshaled, err := bson.Marshal(from)
		if err != nil {
			panic(err)
		}

		err = bson.Unmarshal(marshaled, &jsonFields)
		if err != nil {
			panic(err)
		}

		flattenNestMap("", jsonFields, resultFields, kk)

		return resultFields
	}
}

func fillCamelCaseFields(keys map[string]interface{}, kk map[string]string) {
	for k, v := range keys {
		switch child := v.(type) {
		case map[string]interface{}:
			kk[strings.ToLower(k)] = strcase.ToLowerCamel(k)
			fillCamelCaseFields(child, kk)
		default:
			kk[strings.ToLower(k)] = strcase.ToLowerCamel(k)
		}
	}
}

func flattenNestMap(prefix string, src map[string]interface{}, dest map[string]interface{}, keys map[string]string) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flattenNestMap(prefix+keys[strings.ToLower(k)], child, dest, keys)
		case nil:
			break
		default:
			if k != "id" && k != "_id" {
				dest[prefix+keys[strings.ToLower(k)]] = v
			}
		}
	}
}
