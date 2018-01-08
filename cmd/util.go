package cmd

import (
	"os"
	"github.com/fatih/color"
)

func errorExit(err error) {
	color.Red("Error: %v", err.Error())
	os.Exit(-1)
}
