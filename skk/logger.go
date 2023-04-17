package skk

import (
	"fmt"
	"github.com/fatih/color"
)

func MagentaOnly(format string) {
	color.Magenta(format)
}

func RedOnly(format string) {
	color.Red(format)
}

func Red(f string, e string) {
	fmt.Print(f, ": ")
	color.Red(e)
}

func Blue(f string, e string) {
	fmt.Print(f, ": ")
	color.Blue(e)
}

func Magenta(f string, e string) {
	fmt.Print(f, ": ")
	color.Magenta(e)
}

func Yellow(f string, e string) {
	fmt.Print(f, ": ")
	color.Yellow(e)
}

func Green(f string, e string) {
	fmt.Print(f, ": ")
	color.Green(e)
}
