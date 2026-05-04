package review

type Severity string

const (
	SeverityBlocking    Severity = "blocking"
	SeverityNonBlocking Severity = "non_blocking"
)

type Finding struct {
	ID       string
	Severity Severity
	Summary  string
}
