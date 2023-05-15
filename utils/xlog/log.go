package xlog

import "github.com/fatih/color"

var (
	Debug = color.New(color.FgWhite).SprintFunc()
	Info  = color.New(color.FgGreen).SprintFunc()
	Warn  = color.New(color.FgYellow).SprintFunc()
	Error = color.New(color.FgRed).SprintFunc()
	Fatal = color.New(color.FgRed).SprintFunc()
)
