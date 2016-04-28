package mysql

import (
	"github.com/hpcloud/sidecar-extensions/go/csm"
	"github.com/hpcloud/sidecar-extensions/go/csm/status"
	"github.com/hpcloud/sidecar-extensions/go/extension"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/provisioner"
	"github.com/hpcloud/sidecar-extensions/go/util"
	"github.com/pivotal-golang/lager"
)

const userSize = 16

type mysqlExtension struct {
	prov   provisioner.MySQLProvisioner
	conf   config.MySQLConfig
	logger lager.Logger
}

func NewMySQLExtension(prov provisioner.MySQLProvisioner,
	conf config.MySQLConfig, logger lager.Logger) extension.Extension {
	return &mysqlExtension{prov: prov, conf: conf, logger: logger}
}

func (e *mysqlExtension) CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	dbName := util.NormalizeID(workspaceID)

	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	password, err := util.SecureRandomString(32)
	if err != nil {
		return nil, err
	}

	err = e.prov.CreateUser(dbName, username, password)

	if err != nil {
		return nil, err
	}

	binding := config.MySQLBinding{
		Hostname: e.conf.Host,
		Host:     e.conf.Host,
		Port:     e.conf.Port,
		Username: username,
		User:     username,
		Password: password,
		Database: dbName,
	}

	response := csm.NewCSMResponse(200, binding, status.Successful)
	return &response, err
}
func (e *mysqlExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})
	dbName := util.NormalizeID(workspaceID)
	err := e.prov.CreateDatabase(dbName)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *mysqlExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	err = e.prov.DeleteUser(username)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *mysqlExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	database := util.NormalizeID(workspaceID)
	err := e.prov.DeleteDatabase(database)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *mysqlExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	exists, err := e.prov.IsUserCreated(username)
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
func (e *mysqlExtension) GetWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	database := util.NormalizeID(workspaceID)

	exists, err := e.prov.IsDatabaseCreated(database)
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
