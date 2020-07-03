package model

type App struct {
	Command
	Version string
}

type Command struct {
	Name      string
	Aliases   []string
	Spec      string
	Desc      string
	LongDesc  string
	Options   []Option
	Arguments []Argument
	Commands  []Command
}

type Option struct {
	ShortNames   []string
	LongNames    []string
	Desc         string
	EnvVar       string
	HideValue    bool
	DefaultValue string
}

type Argument struct {
	Name         string
	Desc         string
	EnvVar       string
	HideValue    bool
	DefaultValue string
}
