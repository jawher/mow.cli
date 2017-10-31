package matcher

type shortcut bool

const (
	theShortcut = shortcut(true)
)

func NewShortcut() Matcher {
	return theShortcut
}
func (shortcut) Match(args []string, c *ParseContext) (bool, []string) {
	return true, args
}

func (shortcut) Priority() int {
	return 10
}

func (shortcut) String() string {
	return "*"
}
