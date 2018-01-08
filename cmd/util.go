package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func errorExit(err error) {
	color.Red(err.Error())
	os.Exit(-1)
}

func retrieveVersionFromString(version []byte) (float64, error) {
	fullVersionWithoutQuotes := string(version)
	fullVersion := strings.Split(fullVersionWithoutQuotes, ".")
	majorMinorVersionString := fmt.Sprintf("%s.%s", fullVersion[0], fullVersion[1])

	majorMinorVersion, err := strconv.ParseFloat(majorMinorVersionString, 32)
	if err != nil {
		return 0, fmt.Errorf("error: could not parse docker version: %s", err.Error())
	}
	return majorMinorVersion, nil
}
