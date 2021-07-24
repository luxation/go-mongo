package mongo

import (
	"encoding/json"
	"strconv"
)

func FlattenedMapFromInterface(from interface{}) map[string]interface{} {
	jsonFields := make(map[string]interface{})
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

func flattenNestMap(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			flattenNestMap(prefix+k, child, dest)
		case []interface{}:
			for i := 0; i < len(child); i++ {
				dest[prefix+k+"."+strconv.Itoa(i)] = child[i]
			}
		case nil:
			break
		case string:
			if v != "" {
				dest[prefix+k] = v
			}
		default:
			dest[prefix+k] = v
		}
	}
}
