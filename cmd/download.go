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
	"os/exec"
	"fmt"
	"errors"

	"github.com/spf13/cobra"
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

		if err := dockerInstalled(); err != nil {
			errorExit(err)
		}

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
	dockerCMD := exec.Command("docker")
	if err := dockerCMD.Run(); err != nil {
		return fmt.Errorf("Error running docker, please make sure docker is installed: %s", err.Error())
	}

	dockerPsCMD := exec.Command("docker", "ps")
	if err := dockerPsCMD.Run(); err != nil {
		// Docker error string whenever you run docker ps
		dockerErrorString := "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"
		return errors.New(dockerErrorString)
	}

	return nil
}
