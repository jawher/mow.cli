package matcher

type Matcher interface {
	Match(args []string, c *ParseContext) (bool, []string)
	Priority() int
}

func IsShortcut(matcher Matcher) bool {
	_, ok := matcher.(shortcut)
	return ok
}
