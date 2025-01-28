package log

import "fmt"

type Color int

const (
	DEFAULT Color = iota
	CYAN
	RED
)

const TermColor_Off string = "\033[0m"
const TermCyan string = "\033[0;36m"
const TermRed string = "\033[0;31m"

func Info(msg string) {
	internalPrint(msg, CYAN)
}

func Error(msg string) {
	internalPrint(msg, RED)
}

func internalPrint(msg string, color Color) {
	prefix := getTerminalColor(color)
	suffix := TermColor_Off
	str := prefix + msg + suffix
	fmt.Println(str)
}

func getTerminalColor(color Color) string {
	switch color {
	case CYAN:
		return TermCyan
	case RED:
		return TermRed
	case DEFAULT:
		fallthrough
	default:
		return TermColor_Off
	}
}
