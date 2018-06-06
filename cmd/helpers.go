// TODO: @ganga this file needs a better name or refactors. "helpers ¯\_(ツ)_/¯"
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

func errorExit(err error) {
	color.Red(err.Error())
	os.Exit(-1)
}

func getMachineHardwareName() (string, error) {
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return "", nil
	}

	trimmed := strings.Trim(string(out), "\n")
	return trimmed, nil
}

// ref: https://gist.github.com/indraniel/1a91458984179ab4cf80#gistcomment-2122149
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
			return fmt.Errorf("error extracting tar file: unknown type %s in %s", string(header.Typeflag), header.Name)
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
