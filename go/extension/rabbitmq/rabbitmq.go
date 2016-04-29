package rabbitmq

import (
	"fmt"

	"github.com/hpcloud/sidecar-extensions/go/csm"
	"github.com/hpcloud/sidecar-extensions/go/csm/status"
	"github.com/hpcloud/sidecar-extensions/go/extension"
	"github.com/hpcloud/sidecar-extensions/go/extension/rabbitmq/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/rabbitmq/provisioner"
	"github.com/hpcloud/sidecar-extensions/go/util"
	"github.com/pivotal-golang/lager"
)

const userSize = 16

type rabbitmqExtension struct {
	conf   config.RabbitmqConfig
	logger lager.Logger
	prov   provisioner.RabbitmqProvisionerInterface
}

func NewRabbitmqExtension(prov provisioner.RabbitmqProvisionerInterface, conf config.RabbitmqConfig, logger lager.Logger) extension.Extension {
	return &rabbitmqExtension{prov: prov, conf: conf, logger: logger}
}

func (e *rabbitmqExtension) CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	dbName := util.NormalizeGuid(workspaceID)

	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	password, err := util.SecureRandomString(32)
	if err != nil {
		return nil, err
	}

	credentials, err := e.prov.CreateUser(dbName, username, password)
	if err != nil {
		return nil, err
	}

	binding := config.RabbitmqBinding{
		Hostname:     credentials["host"],
		Host:         credentials["host"],
		VHost:        credentials["vhost"],
		Port:         credentials["port"],
		Username:     credentials["user"],
		Password:     credentials["password"],
		Uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", credentials["user"], credentials["password"], credentials["host"], credentials["port"], credentials["vhost"]),
		DashboardUrl: fmt.Sprintf("http://%s:%s", credentials["host"], credentials["mgmt_port"]),
		Name:         workspaceID,
		User:         credentials["user"],
		Pass:         credentials["password"],
	}

	response := csm.NewCSMResponse(200, binding, status.Successful)
	return &response, err
}
func (e *rabbitmqExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})
	dbName := util.NormalizeGuid(workspaceID)
	err := e.prov.CreateContainer(dbName)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *rabbitmqExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	err = e.prov.DeleteUser(workspaceID, username)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *rabbitmqExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	database := util.NormalizeGuid(workspaceID)
	err := e.prov.DeleteContainer(database)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *rabbitmqExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	exists, err := e.prov.UserExists(workspaceID, username)
	if err != nil {
		return nil, err
	}

	response := csm.CSMResponse{}

	if exists {
		response = csm.NewCSMResponse(200, "", status.Successful)
	} else {
		response = csm.NewCSMResponse(404, "", status.Successful)
	}

	return &response, nil
}
func (e *rabbitmqExtension) GetWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	database := util.NormalizeGuid(workspaceID)

	exists, err := e.prov.ContainerExists(database)
	if err != nil {
		return nil, err
	}

	response := csm.CSMResponse{}

	if exists {
		response = csm.NewCSMResponse(200, "", status.Successful)
	} else {
		response = csm.NewCSMResponse(404, "", status.Successful)
	}

	return &response, nil
}
