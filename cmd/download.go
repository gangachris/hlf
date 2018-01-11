// Copyright © 2018 Chris Ganga <ganga.chris@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	// MinimumDockerVersion is the minimum required docker version for hyperledger fabric to work
	MinimumDockerVersion = 17.03

	// MinimumDockerComposeVersion is the minimum required docker-compose version for hyperledger fabric to work
	MinimumDockerComposeVersion = 1.8

	// FabricVersion is the current stable version of hyperledger fabric
	FabricVersion = "1.0.5"

	// PlatormBinariesURL is the root url for the platform binaries
	PlatormBinariesURL = "https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Hyperledger Fabric Tools (Docker Images, Crypto Tools)",
	Long: `This will download all the required Docker Images to set up a Hyperledger Fabric Environment.
The following Images are downloaded:

	1. hyperledger/fabric-ca
	2. hyperledger/fabric-couchdb
	3. hyperledger/fabric-orderer
	4. hyperledger/fabric-peer
	5. hyperledger/fabric-ccenv
	6. hyperledger/fabric-baseos

The following tools are downloaded:

	1. Cryptogen Tools`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current architecture
		// Check if docker is installed
		// Download Tools
		fmt.Println("download called")

		// check if docker is installed
		if err := dockerInstalled(); err != nil {
			errorExit(err)
		}

		// get system architecture
		// arch := runtime.GOOS + "-" + runtime.GOARCH
		// platformBinariesURL := fmt.Sprintf("%s/%s-%s/hyperledger-fabric-%s-%s.tar.gz", PlatormBinariesURL, arch, FabricVersion, arch, FabricVersion)

		// // // download Platform Binaries
		// if err := downloadPlatformBinaries(platformBinariesURL); err != nil {
		// 	errorExit(err)
		// }

		machineHardwareName, err := getMachineHarwareName()
		if err != nil {
			errorExit(err)
		}

		dockerTag := machineHardwareName + "-" + FabricVersion

		// downloadDockerImages
		if err := downloadDockerImages(dockerTag); err != nil {
			errorExit(err)
		}

		// download correct docker images

		// download cryptoconfig/cryptogen tools

		// TODO: Go should be installed. Maybe serve this as a warning
		// TODO: NodeJS is also a prerequisite, warning maybe
		// TODO: Leave a message that windows is not currently supoorted, but we should try and install windows-build-tools
		// according to the docs

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dockerInstalled() error {
	// check if docker is installed
	dockerCMD := exec.Command("docker")
	if err := dockerCMD.Run(); err != nil {
		return fmt.Errorf("Error running docker, please make sure docker is installed: %s", err.Error())
	}

	// check if docker daemon is running
	dockerPsCMD := exec.Command("docker", "ps")
	if err := dockerPsCMD.Run(); err != nil {
		// Docker error string whenever you run docker ps
		dockerErrorString := "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"
		return errors.New(dockerErrorString)
	}

	// check docker version
	dockerVersionCmdOutput, err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").Output()
	if err != nil {
		return fmt.Errorf("error checking docker version: %s", err.Error())
	}

	dockerVersion, err := retrieveVersionFromCMDOutput(dockerVersionCmdOutput)
	if err != nil {
		return fmt.Errorf("error parsing docker version: %s", err.Error())
	}

	requiredDockerVersion, err := biggerThanMinimumVersion(MinimumDockerVersion, dockerVersion)
	if err != nil {
		return fmt.Errorf("error checking docker version: %s", err.Error())
	}

	if !requiredDockerVersion {
		return fmt.Errorf("error: docker version %.2f.0-ce or higher is required", MinimumDockerVersion)
	}

	// check if docker-compose is installed
	dockerComposeCMDOutput, err := exec.Command("docker-compose", "version", "--short").Output()
	if err != nil {
		return fmt.Errorf("error: please make sure docker-compose is installed: %s", err.Error())
	}

	// check docker-compose version
	dockerComposeVersion, err := retrieveVersionFromCMDOutput(dockerComposeCMDOutput)
	if err != nil {
		return fmt.Errorf("error parsing docker-compose version: %s", err.Error())
	}

	requiredDockerComposeVersion, err := biggerThanMinimumVersion(MinimumDockerComposeVersion, dockerComposeVersion)
	if err != nil {
		return fmt.Errorf("error checking docker-compose version: %s", err.Error())
	}

	if !requiredDockerComposeVersion {
		return fmt.Errorf("error: docker-compose version %.2f.0 or higher is required", MinimumDockerComposeVersion)
	}

	return nil
}

func downloadPlatformBinaries(platormBinariesURL string) error {
	color.Blue("Downloading platform binaries...")
	res, err := http.Get(platormBinariesURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return extractTarGz(res.Body)
}

func downloadDockerImages(dockerTag string) error {
	fabricDockerImages := []string{"peer", "orderer", "couchdb", "ccenv", "javaenv", "kafka", "zookeeper", "tools", "ca"}

	for _, image := range fabricDockerImages {
		imageString := fmt.Sprintf("hyperledger/fabric-%s:%s", image, dockerTag)
		if err := pullDockerImage(imageString); err != nil {
			return err
		}
	}

	color.Green("Docker images downloaded successfully")
	return nil
}

func pullDockerImage(image string) error {
	fmt.Println()
	color.Green("Downloading " + image)
	cmd := exec.Command("docker", "pull", image)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	in := bufio.NewScanner(stdout)

	for in.Scan() {
		fmt.Println(in.Text())
	}

	return in.Err()
}
