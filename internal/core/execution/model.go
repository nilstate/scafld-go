package execution

import "time"

type Request struct {
	Command string
	CWD     string
	Env     []string
	Timeout time.Duration
}

type Result struct {
	ExitCode       int
	Output         string
	DiagnosticPath string
	TimedOut       bool
}
