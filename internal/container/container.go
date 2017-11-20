package container

import "flag"

// ValidationType is the type of validation supported for this field
type ValidationType int

const (
	// Enum is for enum validation, espected to set User/Value/Help
	Enum ValidationType = iota

	// TODO:
	// IPv4Address
	// IPv6Address
	// IPv4Network
	// IPv6Network
	// IntRange
)

/*
Validation contains information needed to validate the value provided by the
user
*/
type Validator struct {
	Type ValidationType

	// TODO:
	// IntMin int
	// IntMax max
	// ...

	// Enums
	// User is the value the user can pass
	User string
	// Value is the value assigned to the string for this user value
	Value string
	// Help is the help message for this value
	Help string
	// Callback is called if the user passes this value to the enum
	Callback func() error
}

/*
Container holds an option or an arg data
*/
type Container struct {
	Name            string
	Desc            string
	EnvVar          string
	Names           []string
	HideValue       bool
	ValueSetFromEnv bool
	ValueSetByUser  *bool
	Value           flag.Value
	DefaultValue    string
	DefaultDisplay  bool
	Validation      []Validator
}
