package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	dockerFilters "github.com/docker/docker/api/types/filters"
	dockerClient "github.com/docker/docker/client"
)

func (tr *TestRunner) findContainer(cli *dockerClient.Client, serviceName string) (string, error) {
	filter := dockerFilters.NewArgs()
	filter.Add("name", serviceName)
	sessiondContainerFilter := dockerTypes.ContainerListOptions{Filters: filter}
	containers, err := cli.ContainerList(context.Background(), sessiondContainerFilter)
	if err != nil {
		return "", err
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("container %s not found ", serviceName)
	}
	return containers[0].ID, nil
}

//RestartService adds ability to restart a particular service managed by docker
func (tr *TestRunner) RestartService(serviceName string) error {
	fmt.Printf("Restarting docker container: %v\n", serviceName)
	ctx := context.Background()
	cli, err := dockerClient.NewEnvClient()
	if err != nil {
		fmt.Printf("error %v getting a new client \n", err)
		return err
	}
	containerID, err := tr.findContainer(cli, serviceName)
	if err != nil {
		fmt.Printf("error %v getting container id \n", err)
		return err
	}
	timeout := 30 * time.Second
	err = cli.ContainerRestart(ctx, containerID, &timeout)
	return err
}

//StopService adds ability to stop a particular service managed by docker
func (tr *TestRunner) PauseService(serviceName string) error {
	fmt.Printf("Pausing docker container: %v\n", serviceName)
	ctx := context.Background()
	cli, err := dockerClient.NewEnvClient()
	if err != nil {
		fmt.Printf("error %v getting a new client \n", err)
		return err
	}
	containerID, err := tr.findContainer(cli, serviceName)
	if err != nil {
		fmt.Printf("error %v getting container id \n", err)
		return err
	}
	err = cli.ContainerPause(ctx, containerID)
	return err
}

//ScanContainerLogs provides ability to scan the container logs for a string
func (tr *TestRunner) ScanContainerLogs(serviceName string, line string) int {
	ctx := context.Background()
	cli, err := dockerClient.NewEnvClient()
	if err != nil {
		fmt.Printf("error %v getting a new client \n", err)
		return 0
	}
	containerID, err := tr.findContainer(cli, serviceName)
	if err != nil {
		fmt.Printf("error %v getting container id \n", err)
		return 0
	}
	reader, _ := cli.ContainerLogs(ctx, containerID, dockerTypes.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	defer reader.Close()

	b, _ := ioutil.ReadAll(reader)
	content := string(b)
	return strings.Count(content, line)
}
