package docker

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/gangachris/hlf/semver"
)

const (
	// MinimumDockerVersion is the minimum required docker version for hyperledger fabric to work
	MinimumDockerVersion = "17.06.2-ce"

	// MinimumDockerComposeVersion is the minimum required docker-compose version for hyperledger fabric to work
	MinimumDockerComposeVersion = "1.14.0"
)

// Client represents a docker client.
type Client struct {
	client *client.Client
}

// New creates an instance of our docker client.
func New() (*Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	c := Client{
		client: cli,
	}

	return &c, nil
}

// Installed checks if Docker is installed
func Installed() error {
	// check if docker is installed
	dockerCMD := exec.Command("docker")
	if err := dockerCMD.Run(); err != nil {
		return fmt.Errorf("error running docker, please make sure docker is installed: %s", err.Error())
	}

	// check if docker daemon is running
	dockerPsCMD := exec.Command("docker", "ps")
	if err := dockerPsCMD.Run(); err != nil {
		// Docker error string whenever you run docker ps
		dockerErrorString := "cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"
		return errors.New(dockerErrorString)
	}

	// check docker version
	dockerVersionCmdOutput, err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").Output()
	if err != nil {
		return fmt.Errorf("error checking docker version: %s", err.Error())
	}

	dockerSemverCeString := strings.TrimSpace(string(dockerVersionCmdOutput))
	// get semver // for now we get rid of -ce
	dockerSemver := dockerSemverCeString[:len(dockerSemverCeString)-3]
	minimumSemver := MinimumDockerVersion[:len(MinimumDockerVersion)-3]

	requiredDockerVersion, err := semver.CorrectVersion(minimumSemver, dockerSemver)
	if err != nil {
		return fmt.Errorf("error checking docker version: %s", err.Error())
	}

	if !requiredDockerVersion {
		return fmt.Errorf("error: docker version %s-ce or higher is required", minimumSemver)
	}

	// check if docker-compose is installed
	dockerComposeCMDOutput, err := exec.Command("docker-compose", "version", "--short").Output()
	if err != nil {
		return fmt.Errorf("error: please make sure docker-compose is installed: %s", err.Error())
	}

	dockerComposeSemverString := strings.TrimSpace(string(dockerComposeCMDOutput))

	// check docker-compose version
	requiredDockerComposeVersion, err := semver.CorrectVersion(MinimumDockerComposeVersion, dockerComposeSemverString)
	if err != nil {
		return fmt.Errorf("error checking docker-compose version: %s", err.Error())
	}

	if !requiredDockerComposeVersion {
		return fmt.Errorf("error: docker-compose version %s or higher is required", MinimumDockerComposeVersion)
	}

	return nil
}

// DownloadDockerImages downloads docker images given a list of images and the tag.
func (c Client) DownloadDockerImages(images []string, tag string) error {
	for _, image := range images {
		color.Blue("Pulling %s", image)
		if err := c.PullAndTagHyperledgerImage(image, tag); err != nil {
			return err
		}
		color.Green("Successfully pulled and tagged %s", image)
	}
	return nil
}

// PullAndTagHyperledgerImage pulls a docker hyperledger image given the name and the tag.
// and tags it with 'latest' tag
func (c Client) PullAndTagHyperledgerImage(imageName, tag string) error {
	imageString := fmt.Sprintf("hyperledger/fabric-%s:%s", imageName, tag)

	_, err := c.client.ImagePull(context.Background(), imageString, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	image := fmt.Sprintf("hyperledger/fabric-%s", imageName)

	return c.client.ImageTag(context.Background(), imageString, image)
}
