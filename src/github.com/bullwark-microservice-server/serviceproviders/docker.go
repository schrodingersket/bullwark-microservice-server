package serviceproviders

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/phayes/freeport"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/dghubble/sling"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/bullwark-microservice-server/common"
	"github.com/bullwark-microservice-server/common/cli"
)

type DockerServiceProvider struct{
	dockerClient *client.Client
}

func NewDockerServiceProvider() DockerServiceProvider {

	dockerClient, err := client.NewClientWithOpts(client.WithVersion("1.39"))

	if err != nil {
		fmt.Printf("WARNING: error creating Docker Client: %s\n", err)
	}

	return DockerServiceProvider{
		dockerClient: dockerClient,
	}
}

func (r DockerServiceProvider) Create(request Request) error {


	extPort, err := freeport.GetFreePort()

	if err != nil {
		return err
	}

	if err := r.createDockerContainer(request, extPort); err != nil {
		return err
	}

	return registerService(request, extPort)
}

func registerService(request Request, extPort int) error {

	var config = common.Configs[cli.RegistrarConfigType].(cli.RegistrarConfig)
	request.Port = extPort

	// Generate registration URL
	//
	var registrarUrl = url.URL{
		Scheme: *config.Scheme,
		Host:   fmt.Sprintf("%s:%d", *config.Host, *config.Port),
		Path:   fmt.Sprintf("%s/register", *config.BaseURL),
	}

	// Send request
	//
	var responseBody interface{}
	req, err := sling.New().
		Post(registrarUrl.String()).
		BodyJSON(request).
		ReceiveSuccess(responseBody)

	if err != nil {
		return err
	}

	if req.StatusCode/100 != 2 {
		return errors.New(fmt.Sprintf("Request failed with status code"+
			" %d.", req.StatusCode))
	} else {
		fmt.Println("Success.")
	}

	return nil
}

func (r DockerServiceProvider) createDockerContainer(request Request, extPort int) error {

	var coreConfig = common.Configs[cli.CoreConfigType].(cli.CoreConfig)
	var internalJarPath = fmt.Sprintf("/var/run/%s.jar", request.ServiceKey)
	var bgContext = context.Background()

	// Pull image
	//
	out, err := r.dockerClient.ImagePull(bgContext, "java:8", types.ImagePullOptions{})

	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, out)

	if err != nil {
		return err
	}

	var config = container.Config{
		Image: "java:8",
		Cmd: []string{
			"java",
			"-jar",
			internalJarPath,
		},
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}

	fileExt := filepath.Ext(request.Filepath)

	portMap := nat.PortMap{
		"8080/tcp": []nat.PortBinding{
			{
				HostIP: "0.0.0.0",
				HostPort: fmt.Sprintf("%d", extPort),
			},
		},
	}

	var hostConfig = container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s/%s/%s:%s",
				*coreConfig.BinSavePath,
				fileExt,
				request.Filepath,
				internalJarPath),
		},
		PublishAllPorts: true,
		PortBindings: portMap,
	}

	var networkingConfig = network.NetworkingConfig{
	}

	c, err := r.dockerClient.ContainerCreate(context.Background(),
		&config,
		&hostConfig,
		&networkingConfig,
		"")

	if err != nil {
		return err
	}

	return r.dockerClient.ContainerStart(context.Background(),
		c.ID,
		types.ContainerStartOptions{})
}
