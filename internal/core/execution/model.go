package execution

import "time"

type Request struct {
	Command        string
	CWD            string
	Env            []string
	Timeout        time.Duration
	IdleTimeout    time.Duration
	TerminateGrace time.Duration
}

type Result struct {
	ExitCode       int
	Output         string
	Stdout         string
	Stderr         string
	DiagnosticPath string
	TimedOut       bool
	KillReason     string
}
