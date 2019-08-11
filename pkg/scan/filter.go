package scan

type ResultFilter interface {
	ShouldIgnore(Result) bool
}
