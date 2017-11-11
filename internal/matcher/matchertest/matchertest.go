package matchertest

import (
	"strings"

	"fmt"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/jawher/mow.cli/internal/matcher"
	"github.com/jawher/mow.cli/internal/values"
)

// NewArg creates a positional argument matcher given its name, e.g. SRC
func NewArg(name string) matcher.Matcher {
	var s string
	con := &container.Container{
		Name:  name,
		Names: []string{name},
		Value: values.NewString(&s, ""),
	}

	return matcher.NewArg(con)
}

// NewOpt creates a short and long option matcher given its (space separated) names, e.g. -f --force
func NewOpt(name string) matcher.Matcher {
	s := ""
	names := strings.Fields(name)
	con := &container.Container{
		Name:  name,
		Names: names,
		Value: values.NewString(&s, ""),
	}

	index := map[string]*container.Container{}
	for _, n := range names {
		index[n] = con
	}

	return matcher.NewOpt(con, index)
}

// NewOptions create an options matcher given their names, e.g. -abc
func NewOptions(names string) matcher.Matcher {
	names = strings.TrimPrefix(names, "-")

	var cons []*container.Container
	index := map[string]*container.Container{}

	for _, r := range names {
		name := fmt.Sprintf("%c", r)
		var s string
		fname := "-" + name
		con := &container.Container{
			Name:  fname,
			Names: []string{fname},
			Value: values.NewString(&s, ""),
		}
		cons = append(cons, con)
		index[fname] = con

	}

	return matcher.NewOptions(cons, index)
}
