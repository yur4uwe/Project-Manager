package parse_json

import (
	"strings"
)

func Parse(JSON string) map[string]string {
	var kvps = make(map[string]string)

	JSON = strings.Trim(JSON, "[]")
	JSON = strings.Trim(JSON, "{}")

	var pairs = strings.Split(JSON, ",")

	for _, pair := range pairs {
		var key, value = ParseProperty(pair)

		kvps[key] = value
	}

	return kvps
}

func ParseProperty(property string) (string, string) {
	var split_property = strings.Split(property, ":")
	var key = split_property[0]
	var value = strings.Join(split_property[1:], ":")

	key = strings.Trim(key, " ")
	value = strings.Trim(value, " ")

	key = strings.Replace(key, "\"", "", -1)
	value = strings.Replace(value, "\"", "", -1)

	return key, value
}
