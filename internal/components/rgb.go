package components

import "fmt"

// rgb converts RGB components to a truecolor ANSI prefix string.
func rgb(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}
