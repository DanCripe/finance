package main

const (
	reset  = "\033[0m"
	red    = "\033[31;1m"
	yellow = "\033[33;1m"
	green  = "\033[32;1m"
	blue   = "\033[34;1m"
)

func colorize(val float64, good, bad string) string {
	if val == 0.0 {
		return bad + p.Sprintf("%13.2f", val) + reset
	}
	return good + p.Sprintf("%13.2f", val) + reset
}
