package main

import (
	"fmt"
	"os"
	"strings"
)

func handleErr(errMsg string) {
	// show cursor
	fmt.Fprint(os.Stdout, "\x1b[?25h")
	fmt.Fprintf(os.Stderr, "%s %s\n", Bold(Red("error:")), errMsg)
	os.Exit(1)
}

func Black(text ...string) string {
	return fmt.Sprintf("\x1b[30m%s\x1b[0m", strings.Join(text, " "))
}

func Red(text ...string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", strings.Join(text, " "))
}

func Green(text ...string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", strings.Join(text, " "))
}

func Yellow(text ...string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", strings.Join(text, " "))
}

func Blue(text ...string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", strings.Join(text, " "))
}

func Purple(text ...string) string {
	return fmt.Sprintf("\x1b[35m%s\x1b[0m", strings.Join(text, " "))
}

func Cyan(text ...string) string {
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", strings.Join(text, " "))
}

func Gray(text ...string) string {
	return fmt.Sprintf("\x1b[37m%s\x1b[0m", strings.Join(text, " "))
}

func Italic(text ...string) string {
	return "\x1b[3m" + strings.Join(text, " \x1b[3m") + "\x1b[0m"
}

func Bold(text ...string) string {
	return "\x1b[1m" + strings.Join(text, " \x1b[1m") + "\x1b[0m"
}

func Unwrap[T any](val T, err error) T {
	if err != nil { handleErr(err.Error()) }

	return val
}

func UnwrapOr[T any](val T, err error) func(T) T {
	if err != nil {
		return func(d T) T {
			return d
		}
	} else {
		return func(_ T) T {
			return val
		}
	}
}

func UnwrapOrElse[T any](val T, err error) func(func() T) T {
	if err != nil {
		return func(fn func() T) T {
			return fn()
		}
	} else {
		return func(_ func() T) T {
			return val
		}
	}

}

func Expect(err error) {
	if err != nil { handleErr(err.Error()) }
}
