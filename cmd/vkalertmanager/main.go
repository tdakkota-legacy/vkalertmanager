package main

import (
	"os"
)

func main() {
	if err := NewApp().Run(os.Args); err != nil {
		_, _ = os.Stdout.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
