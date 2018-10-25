package log

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

//Log checks and logs a error
func LogError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

// print something to stdout at normal verbosity
func Info(message string) {
	fmt.Println(message)
}

func Infof(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", a...)
}

// for now just print it
func Verbose(message string) {
	fmt.Println(message)
}

func FatalExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func FatalPanic(err error) {
	if err != nil {
		panic(errors.Wrap(err, ""))
	}
}
