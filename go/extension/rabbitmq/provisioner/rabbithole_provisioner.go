package provisioner

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hpcloud/sidecar-extensions/go/extension/rabbitmq/config"
	"github.com/hpcloud/sidecar-extensions/go/util"
	"github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-golang/lager"

	dockerclient "github.com/fsouza/go-dockerclient"
)

const CONTAINER_START_TIMEOUT int = 30

type RabbitHoleProvisioner struct {
	rabbitmqConfig config.RabbitmqConfig
	client         *dockerclient.Client
	logger         lager.Logger
	connected      bool
}

func NewRabbitHoleProvisioner(logger lager.Logger, conf config.RabbitmqConfig) RabbitmqProvisionerInterface {
	return &RabbitHoleProvisioner{logger: logger, rabbitmqConfig: conf}
}

func (provisioner *RabbitHoleProvisioner) CreateContainer(containerName string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	err = provisioner.pullImage(provisioner.rabbitmqConfig.DockerImage, provisioner.rabbitmqConfig.ImageVersion)
	if err != nil {
		return err
	}

	admin_user, err := util.SecureRandomString(32)
	if err != nil {
		return err
	}
	admin_pass, err := util.SecureRandomString(32)
	if err != nil {
		return err
	}
	hostConfig := dockerclient.HostConfig{PublishAllPorts: true}
	createOpts := dockerclient.CreateContainerOptions{
		Config: &dockerclient.Config{
			Image: provisioner.rabbitmqConfig.DockerImage + ":" + provisioner.rabbitmqConfig.ImageVersion,
			Env: []string{"RABBITMQ_DEFAULT_USER=" + admin_user,
				"RABBITMQ_DEFAULT_PASS=" + admin_pass},
		},
		HostConfig: &hostConfig,
		Name:       containerName,
	}

	container, err := provisioner.client.CreateContainer(createOpts)
	if err != nil {
		return err
	}

	provisioner.client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		return err
	}

	retry := 1
	for retry < CONTAINER_START_TIMEOUT {
		state, err := provisioner.getContainerState(containerName)
		if err != nil {
			return err
		}
		if state.Running {
			break
		}
		retry++
	}

	return nil
}

func (provisioner *RabbitHoleProvisioner) DeleteContainer(containerName string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	containerID, err := provisioner.getContainerId(containerName)
	if err != nil {
		return err
	}

	err = provisioner.client.StopContainer(containerID, 5)
	if err != nil {
		return err
	}

	return provisioner.client.RemoveContainer(dockerclient.RemoveContainerOptions{
		ID:    containerID,
		Force: true,
	})
}

func (provisioner *RabbitHoleProvisioner) ContainerExists(containerName string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}

	_, err = provisioner.getContainer(containerName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (provisioner *RabbitHoleProvisioner) DeleteUser(containerName, user string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	host, err := provisioner.getHost()
	if err != nil {
		return err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return err
	}

	_, err = rmqc.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (provisioner *RabbitHoleProvisioner) UserExists(containerName, user string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}

	host, err := provisioner.getHost()
	if err != nil {
		return false, err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return false, err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return false, err
	}

	users, err := rmqc.ListUsers()
	if err != nil {
		return false, err
	}
	if users == nil {
		return false, err
	}

	for _, u := range users {
		if u.Name == user {
			return true, nil
		}
	}

	return false, nil
}

func (provisioner *RabbitHoleProvisioner) getClient() (*dockerclient.Client, error) {
	client, err := dockerclient.NewClient(provisioner.rabbitmqConfig.DockerEndpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (provisioner *RabbitHoleProvisioner) pullImage(imageName, version string) error {
	var buf bytes.Buffer
	pullOpts := dockerclient.PullImageOptions{
		Repository:   imageName,
		Tag:          version,
		OutputStream: &buf,
	}

	err := provisioner.client.PullImage(pullOpts, dockerclient.AuthConfiguration{})
	if err != nil {
		return err
	}
	return nil
}

func (provisioner *RabbitHoleProvisioner) findImage(imageName string) (*dockerclient.Image, error) {
	image, err := provisioner.client.InspectImage(imageName)
	if err != nil {
		return nil, fmt.Errorf("Could not find base image %s: %s", imageName, err.Error())
	}

	return image, nil
}

func (provisioner *RabbitHoleProvisioner) getContainerId(containerName string) (string, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return "", err
	}
	return container.ID, nil
}

func (provisioner *RabbitHoleProvisioner) getContainer(containerName string) (dockerclient.APIContainers, error) {
	opts := dockerclient.ListContainersOptions{
		All: true,
	}
	containers, err := provisioner.client.ListContainers(opts)
	if err != nil {
		return dockerclient.APIContainers{}, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if strings.TrimPrefix(n, "/") == containerName {
				return c, nil
			}
		}
	}

	return dockerclient.APIContainers{}, fmt.Errorf("Container %s not found", containerName)
}

func (provisioner *RabbitHoleProvisioner) inspectContainer(containerId string) (*dockerclient.Container, error) {
	return provisioner.client.InspectContainer(containerId)
}

func (provisioner *RabbitHoleProvisioner) getAdminCredentials(containerName string) (map[string]string, error) {

	m := make(map[string]string)
	containerId, err := provisioner.getContainerId(containerName)
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return nil, err
	}

	container, err := provisioner.inspectContainer(containerId)
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return nil, err
	}

	var env dockerclient.Env
	env = make([]string, len(container.Config.Env)) // container.Config.Env.(dockerclient.Env)  // dockerclient.Env( []string{ container.Config.Env })
	copy(env, container.Config.Env)
	m["user"] = env.Get("RABBITMQ_DEFAULT_USER")
	m["password"] = env.Get("RABBITMQ_DEFAULT_PASS")
	for k, v := range container.NetworkSettings.Ports {
		if k == "15672/tcp" {
			m["mgmt_port"] = v[0].HostPort
		}
		if k == "5672/tcp" {
			m["port"] = v[0].HostPort
		}
	}
	return m, nil
}

func (provisioner *RabbitHoleProvisioner) getContainerState(containerName string) (dockerclient.State, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return dockerclient.State{}, nil
	}

	c, err := provisioner.inspectContainer(container.ID)
	if err != nil {
		return dockerclient.State{}, nil
	}
	return c.State, nil
}

func (provisioner *RabbitHoleProvisioner) CreateUser(containerName, newUser, userPass string) (map[string]string, error) {
	host, err := provisioner.getHost()
	if err != nil {
		return nil, err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return nil, err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)

	_, err = rmqc.PutUser(newUser, rabbithole.UserSettings{Password: userPass, Tags: "management,policymaker"})
	if err != nil {
		return nil, err
	}

	_, err = rmqc.UpdatePermissionsIn("/", newUser, rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"})
	if err != nil {
		return nil, err
	}
	m["host"] = host
	m["user"] = newUser
	m["password"] = userPass
	m["mgmt_port"] = admin["mgmt_port"]
	m["port"] = admin["port"]
	x, err := rmqc.GetVhost("/")
	if err != nil {
		return nil, err
	}
	m["vhost"] = x.Name

	return m, nil
}

func (provisioner *RabbitHoleProvisioner) getHost() (string, error) {
	host := ""
	dockerUrl, err := url.Parse(provisioner.rabbitmqConfig.DockerEndpoint)
	if err != nil {
		return "", err
	}

	if dockerUrl.Scheme == "unix" {
		host, err = util.GetLocalIP()
		if err != nil {
			return "", err
		}
	} else {
		host = strings.Split(dockerUrl.Host, ":")[0]
	}

	return host, nil
}

func (provisioner *RabbitHoleProvisioner) connect() error {
	if provisioner.connected {
		return nil
	}

	var err error

	dockerUrl, err := url.Parse(provisioner.rabbitmqConfig.DockerEndpoint)
	if err != nil {
		return err
	}

	if dockerUrl.Scheme == "" {
		return errors.New("Invalid URL format")
	}

	provisioner.client, err = provisioner.getClient()
	if err != nil {
		return err
	}

	provisioner.connected = true
	return nil
}
