package cmd

import (
	"math"
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

func retrieveVersionFromCMDOutput(version []byte) (float64, error) {
	fullVersionWithoutQuotes := string(version)
	fullVersion := strings.Split(fullVersionWithoutQuotes, ".")
	majorMinorVersionString := fmt.Sprintf("%s.%s", fullVersion[0], fullVersion[1])

	majorMinorVersion, err := strconv.ParseFloat(majorMinorVersionString, 32)
	if err != nil {
		return 0, fmt.Errorf("error: could not parse docker version: %s", err.Error())
	}
	return majorMinorVersion, nil
}

func biggerThanMinimumVersion(minimum, current float64) (bool, error) {
	minimumFlooredNumber, currentFlooredNumber := math.Floor(minimum), math.Floor(current)
	if (currentFlooredNumber < minimumFlooredNumber) {
		return false, nil
	}

	// check second digit
	minimumString := fmt.Sprintf("%f", minimum)
	splitMinimumString := strings.Split(minimumString, ".")
	minimumDecimalString := strings.Split(splitMinimumString[1], "00")
	var minimumDecimal int
	var err error
	if len(minimumDecimalString) != 0 {
		minimumDecimal, err = strconv.Atoi(minimumDecimalString[0])
		if err != nil {
			return false, err
		}
	}

	currentString := fmt.Sprintf("%f", current)
	splitCurrentString := strings.Split(currentString, ".")
	currentDecimalString := strings.Split(splitCurrentString[1], "00")
	currentDecimal := 0
	if len(currentDecimalString) != 0 {
		currentDecimal, err = strconv.Atoi(currentDecimalString[0])
		if err != nil {
			return false, err
		}
	}

	return currentDecimal >= minimumDecimal, nil
}
