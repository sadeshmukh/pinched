package secrets

import (
	"regexp"
)

func Substitute(input string) string {
	// {{NAME}} replaces all these thingies

	match := regexp.MustCompile(`\{\{(\w+)\}\}`)

	return match.ReplaceAllStringFunc(
		input, func(secret string) string {
			val, err := GetSecret(secret[2 : len(secret)-2])
			if err != nil {
				return secret
			}
			return val
		},
	)
}
