package tui

import (
	"regexp"
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// MatchFilter applies the live filter grammar.
// Tokens are combined with AND.
// status=Running, ns=production, namespace=production, label=value, /regex/, and fuzzy text are supported.
func MatchFilter(row kubernetes.ResourceRow, input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return true
	}
	for _, token := range strings.Fields(input) {
		if !matchToken(row, token) {
			return false
		}
	}
	return true
}

func matchToken(row kubernetes.ResourceRow, token string) bool {
	if len(token) >= 2 && strings.HasPrefix(token, "/") && strings.HasSuffix(token, "/") {
		pattern, err := regexp.Compile(token[1 : len(token)-1])
		if err != nil {
			return false
		}
		return pattern.MatchString(row.Name + " " + row.Namespace + " " + row.Status)
	}
	key, value, field := strings.Cut(token, "=")
	if field {
		key = strings.ToLower(key)
		switch key {
		case "status":
			return strings.EqualFold(row.Status, value)
		case "ns", "namespace":
			return strings.EqualFold(row.Namespace, value)
		default:
			return strings.EqualFold(row.Labels[key], value)
		}
	}
	return fuzzy(strings.ToLower(row.Name+" "+row.Namespace+" "+row.Status+" "+row.Kind+" "+row.Extra["node"]), strings.ToLower(token))
}

func fuzzy(value, query string) bool {
	if strings.Contains(value, query) {
		return true
	}
	at := 0
	for _, char := range query {
		next := strings.IndexRune(value[at:], char)
		if next < 0 {
			return false
		}
		at += next + 1
	}
	return true
}
