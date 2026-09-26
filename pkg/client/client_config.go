package client

import "fmt"

// IsVerbose will Toggle HTTP output
var IsVerbose bool

func SetVerbosity(verbose bool) {
	IsVerbose = verbose
}

func Verbose(msg string) {
	if IsVerbose {
		_, _ = fmt.Printf("[HTTP]: %s\n", msg)
	}
}

var GlobalRetriesPolicy Retries

func SetGlobalRetries(retries Retries) {
	GlobalRetriesPolicy = retries
}

type Retries struct {
	Timeout       int
	CustomMessage map[string]int
}
