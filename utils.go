package cli

import (
	"flag"
	"os"
	"strings"
)

func setFromEnv(into flag.Value, envVars string) bool {
	multiValued, isMulti := into.(multiValued)

	if len(envVars) > 0 {
		for _, rev := range strings.Split(envVars, " ") {
			ev := strings.TrimSpace(rev)
			if len(ev) == 0 {
				continue
			}

			v := os.Getenv(ev)
			if len(v) == 0 {
				continue
			}
			if !isMulti {
				if err := into.Set(v); err == nil {
					return true
				}
				continue
			}

			vs := strings.Split(v, ",")
			if err := setMultivalued(multiValued, vs); err == nil {
				return true
			}
		}
	}
	return false
}

func setMultivalued(into multiValued, values []string) error {
	into.Clear()

	for _, v := range values {
		v = strings.TrimSpace(v)
		if err := into.Set(v); err != nil {
			into.Clear()
			return err
		}
	}

	return nil
}

func words(v string) []string {
	parts := strings.Split(v, " ")
	res := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		res = append(res, part)
	}
	return res
}
