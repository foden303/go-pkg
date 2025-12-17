package logger

import "strings"

var keyReplacer = strings.NewReplacer(
	".", "_",
	"-", "_",
	"|", "_",
)

func sanitizeKey(key string) string {
	return keyReplacer.Replace(key)
}

func normalizeTags(tags map[string]interface{}) map[string]interface{} {
	newTags := make(map[string]interface{}, len(tags))
	for k, v := range tags {
		key := strings.TrimSpace(strings.ToLower(sanitizeKey(k)))
		newTags[sanitizeKey(key)] = v
	}
	return newTags
}
