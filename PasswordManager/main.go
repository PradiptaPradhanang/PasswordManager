package main

import (
	"os"
	"passmana/appfrontend"
	"passmana/cmd"
)

func main() {
	if len(os.Args) > 1 {
		// If any arguments are passed, assume CLI mode
		cmd.Execute()
	} else {
		// No arguments? Launch GUI
		appfrontend.EntryPoint()
	}
}
