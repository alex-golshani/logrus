package logrus

type formatMode int

const (
	formatted formatMode = iota
	unformatted
	newLine
)

