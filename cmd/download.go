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
	"fmt"
	"net/http"
	"runtime"

	"github.com/fatih/color"
	"github.com/gangachris/hlf-cli/docker"
	"github.com/spf13/cobra"
)

const (
	// FabricVersion is the current stable version of hyperledger fabric
	FabricVersion = "1.1.0"

	// ThirdPartyVersionTag represents the version of third party images released (couchdb, kafka and zookeeper)
	ThirdPartyVersionTag = "0.4.6"

	// PlatformBinariesURL is the root url for the platform binaries
	PlatformBinariesURL = "https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric"

	// HYPERLEDGER string used for docker hub when pulling images
	HYPERLEDGER = "hyperledger"
)

// downloadCmd will download the platform binaries and the docker images
// once the download is done, the images are tagged
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Hyperledger Fabric Tools (Docker Images, Platform Binaries)",
	Long: `This will download all the required Docker Images to set up a Hyperledger Fabric Environment.
The following Images are downloaded:

	 1. hyperledger/fabric-ca
	 2. hyperledger/fabric-couchdb
	 3. hyperledger/fabric-orderer
	 4. hyperledger/fabric-peer
	 5. hyperledger/fabric-ccenv
	 6. hyperledger/fabric-baseos
	 7. hyperledger/fabric-javaenv
	 8. hyperledger/fabric-tools
	 9. hyperledger/fabric-zookeeper
	10. hyperledger/fabric-kafka

The following binaries tools are also downloaded:

	1. configtxgen
	2. configtxlator
	3. cryptogen
	4. peer
	5. orderer`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue("Downloading Platform Binaries and Docker Images")

		// check if docker is installed
		if err := docker.Installed(); err != nil {
			errorExit(err)
		}

		// get system architecture
		arch := runtime.GOOS + "-" + runtime.GOARCH
		platformBinariesURL := fmt.Sprintf("%s/%s-%s/hyperledger-fabric-%s-%s.tar.gz", PlatformBinariesURL, arch, FabricVersion, arch, FabricVersion)

		// download Platform Binaries
		// TODO: @ganga maybe add to path???
		// TODO: @ganga maybe show progress with uilive
		// github.com/gosuri/uilive
		// TODO: @ganga a way to check if platform binaries have been downloaded
		if err := downloadPlatformBinaries(platformBinariesURL); err != nil {
			errorExit(err)
		}

		machineHardwareName, err := getMachineHardwareName()
		if err != nil {
			errorExit(err)
		}

		dockerClient, err := docker.New()
		if err != nil {
			errorExit(err)
		}

		fabricDockerImages := []string{"peer", "orderer", "ccenv", "javaenv", "tools", "ca"}
		fabricTag := machineHardwareName + "-" + FabricVersion

		if err := dockerClient.DownloadDockerImages(fabricDockerImages, fabricTag); err != nil {
			errorExit(err)
		}

		thirdPartyDockerImages := []string{"couchdb", "zookeeper", "kafka"}
		thirdPartyTag := machineHardwareName + "-" + ThirdPartyVersionTag
		if err := dockerClient.DownloadDockerImages(thirdPartyDockerImages, thirdPartyTag); err != nil {
			errorExit(err)
		}

		color.Green("Successfully downloaded docker images")

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

func downloadPlatformBinaries(platformBinariesURL string) error {
	color.Blue("Downloading platform binaries...")
	res, err := http.Get(platformBinariesURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return extractTarGz(res.Body)
}
