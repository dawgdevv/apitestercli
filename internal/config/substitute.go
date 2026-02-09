package config

import (
	"fmt"
	"regexp"
)

var re = regexp.MustCompile(`\{\{(\w+)\}\}`)

// ResolveEnv resolves inter-variable references within an env map.
// e.g. if base_url = "http://localhost" and endpoint = "{{base_url}}/api",
// endpoint resolves to "http://localhost/api".
func ResolveEnv(env map[string]string) map[string]string {
	resolved := make(map[string]string, len(env))
	for k, v := range env {
		resolved[k] = v
	}
	// Multi-pass resolution (max 10 iterations to avoid infinite loops)
	for i := 0; i < 10; i++ {
		changed := false
		for k, v := range resolved {
			newVal := re.ReplaceAllStringFunc(v, func(match string) string {
				key := re.FindStringSubmatch(match)[1]
				if val, ok := resolved[key]; ok && key != k { // avoid self-reference
					return val
				}
				return match
			})
			if newVal != resolved[k] {
				resolved[k] = newVal
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return resolved
}

func SubstituteString(input string, env map[string]string) (string, error) {
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		key := re.FindStringSubmatch(match)[1]

		if val, ok := env[key]; ok {
			return val
		}

		return match
	})

	if re.MatchString(result) {
		return "", fmt.Errorf("unresolved variable in string: %s", input)
	}

	return result, nil

}
