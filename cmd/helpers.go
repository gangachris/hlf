package cmd

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
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

// Check if there are three version numbers major, minor, patch
type semanticVersion struct {
	major, minor, patch int
}

var errInvalidSemver = errors.New("invalid semver")

func getSemverFromString(semver string) (semanticVersion, error) {
	smv := semanticVersion{}

	splitSemver := strings.Split(semver, ".")
	if len(splitSemver) == 0 {
		return smv, errInvalidSemver
	}

	if len(splitSemver) == 3 {
		// good
		return getSemverFromSlice(splitSemver)
	}

	var s []string
	if len(splitSemver) == 2 {
		s = append(s, splitSemver...)
		s = append(s, "0")
		return getSemverFromSlice(s)
	}

	if len(splitSemver) == 1 {
		s = append(s, splitSemver...)
		s = append(s, "0", "0")
		return getSemverFromSlice(s)
	}

	return smv, errInvalidSemver
}

func getSemverFromSlice(semver []string) (semanticVersion, error) {
	smv := semanticVersion{}
	major, err := strconv.Atoi(semver[0])
	if err != nil {
		return smv, err
	}
	smv.major = major

	minor, err := strconv.Atoi(semver[1])
	if err != nil {
		return smv, err
	}
	smv.minor = minor

	patch, err := strconv.Atoi(semver[2])
	if err != nil {
		return smv, err
	}
	smv.patch = patch

	return smv, nil
}

// func correctSemver(minimum, current string) (bool, error) {

// 	var min, max semver

// 	minSplit := strings.Split(minimum, ".")
// 	major, err := strconv.Atoi(minSplit[0])
// 	if err != nil {
// 		return false, err
// 	}
// 	min.major = major

// 	minor, err := strconv.Atoi(minSplit[1])
// 	if err != nil {
// 		return false, err
// 	}
// 	min.minor = minor

// 	patch, err := strconv.Atoi(minSplit[2])
// 	if err != nil {
// 		return false, err
// 	}
// 	min.patch = patch

// 	currSplit := strings.Split(current, ".")
// 	major, err = strconv.Atoi(currSplit[0])
// 	if err != nil {
// 		return false, err
// 	}
// 	max.major = major

// 	minor, err = strconv.Atoi(currSplit[1])
// 	if err != nil {
// 		return false, err
// 	}
// 	max.minor = minor

// 	patch, err = strconv.Atoi(currSplit[2])
// 	if err != nil {
// 		return false, err
// 	}
// 	max.patch = patch
// }

func getMachineHarwareName() (string, error) {
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return "", nil
	}

	trimmed := strings.Trim(string(out), "\n")
	return trimmed, nil
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
			err := os.Mkdir(platformsBinariesDir+"/"+header.Name, 0755)

			if os.IsExist(err) {
				return nil
			}

			if err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(platformsBinariesDir + "/" + header.Name)
			if os.IsExist(err) {
				return nil
			}

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

	err = os.MkdirAll(platformBinariesDir, 0755)

	if os.IsExist(err) {
		return platformBinariesDir, nil
	}

	if err != nil {
		return "", err
	}

	return platformBinariesDir, nil
}
