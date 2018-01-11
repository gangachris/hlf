package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
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
	if currentFlooredNumber < minimumFlooredNumber {
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

func getMachineHarwareName() (string, error) {
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return "", nil
	}

	return string(out), nil
}

// tip: https://gist.github.com/indraniel/1a91458984179ab4cf80#gistcomment-2122149
func extractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		platformsBinariesDir, err := getPlatformBinariesDir()
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(platformsBinariesDir+"/"+header.Name, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(platformsBinariesDir + "/" + header.Name)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("error extracting tar file: unkown type %s in %s", string(header.Typeflag), header.Name)
		}
	}

	return nil
}

func getPlatformBinariesDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	platformBinariesDir := home + "/.hlf-cli"

	if err := os.MkdirAll(platformBinariesDir, 0755); err != nil {
		return "", err
	}

	return platformBinariesDir, nil
}
