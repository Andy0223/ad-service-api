package redis

import (
	"sort"
)

func GenerateRedisKey(key map[string]string) string {
	redisKey := "ads"

	// Create a slice of keys and sort it
	keys := make([]string, 0, len(key))
	for k := range key {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate over keys in sorted order
	for _, k := range keys {
		redisKey += ":" + k + ":" + key[k]
	}

	return redisKey
}
