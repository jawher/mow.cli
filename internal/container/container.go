package container

import "flag"

type Container struct {
	Name            string
	Desc            string
	EnvVar          string
	Names           []string
	HideValue       bool
	ValueSetFromEnv bool
	ValueSetByUser  *bool
	Value           flag.Value
}
