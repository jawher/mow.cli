package matcher

type optsEnd bool

func NewOptsEnd() Matcher {
	return theOptsEnd
}

const (
	theOptsEnd = optsEnd(true)
)

func (u optsEnd) Match(args []string, c *ParseContext) (bool, []string) {
	c.RejectOptions = true
	return true, args
}

func (optsEnd) Priority() int {
	return 9
}

func (optsEnd) String() string {
	return "--"
}
