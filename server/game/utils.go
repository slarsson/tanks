package game

func validateName(str *string) bool {
	return true
}

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)

const (
	colorReset  = string("\033[0m")
	colorRed    = string("\033[1;31m")
	colorGreen  = string("\033[1:32m")
	colorYellow = string("\033[1;33m")
	colorBlue   = string("\033[1;34m")
	colorPurple = string("\033[1;35m")
	colorCyan   = string("\033[1;36m")
	colorWhite  = string("\033[1;37m")
)
