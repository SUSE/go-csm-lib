package postgres

import (
	"github.com/hpcloud/sidecar-extensions/go/csm"
	"github.com/hpcloud/sidecar-extensions/go/csm/status"
	"github.com/hpcloud/sidecar-extensions/go/extension"
	"github.com/hpcloud/sidecar-extensions/go/extension/postgres/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/postgres/provisioner"
	"github.com/hpcloud/sidecar-extensions/go/util"
	"github.com/pivotal-golang/lager"
)

const userSize = 16

type postgresExtension struct {
	prov   provisioner.PostgresProvisionerInterface
	conf   config.PostgresConfig
	logger lager.Logger
}

func NewPostgresExtension(prov provisioner.PostgresProvisionerInterface,
	conf config.PostgresConfig, logger lager.Logger) extension.Extension {
	return &postgresExtension{prov: prov, conf: conf, logger: logger}
}

func (e *postgresExtension) CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
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

	err = e.prov.CreateUser(dbName, username, password)
	if err != nil {
		return nil, err
	}

	binding := config.PostgresBindingCredentials{
		Hostname:         e.conf.Host,
		Host:             e.conf.Host,
		Database:         dbName,
		Password:         password,
		Port:             e.conf.Port,
		Username:         username,
		ConnectionString: config.GenerateConnectionString(config.ConnectionStringTemplate, e.conf.Host, e.conf.Port, dbName, username, password),
		Name:             dbName,
		User:             username,
		Uri:              config.GenerateConnectionString(config.UriTemplate, e.conf.Host, e.conf.Port, dbName, username, password),
		JdbcUrl:          config.GenerateConnectionString(config.JdbcUrilTemplate, e.conf.Host, e.conf.Port, dbName, username, password),
	}

	response := csm.NewCSMResponse(200, binding, status.Successful)
	return &response, err
}

func (e *postgresExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})
	dbName := util.NormalizeGuid(workspaceID)
	err := e.prov.CreateDatabase(dbName)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *postgresExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	dbName := util.NormalizeGuid(workspaceID)
	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	err = e.prov.DeleteUser(dbName, username)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *postgresExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	database := util.NormalizeGuid(workspaceID)
	err := e.prov.DeleteDatabase(database)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *postgresExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	exists, err := e.prov.UserExists(username)
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
func (e *postgresExtension) GetWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	database := util.NormalizeGuid(workspaceID)

	exists, err := e.prov.DatabaseExists(database)
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
