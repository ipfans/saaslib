package liberrors

import "github.com/morikuni/failure"

var (
	ErrConfigNotEnabled failure.StringCode = "ConfigNotEnabled"
	ErrConfigReadFailed failure.StringCode = "ConfigReadFailed"
)
