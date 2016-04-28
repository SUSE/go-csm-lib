package main

import (
	"os"

	"github.com/hpcloud/sidecar-extensions/go/csm"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/provisioner"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("mysql-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.MySQLConfig{}
	env.Parse(&conf)

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Error("main", err)
		os.Exit(1)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath)
	prov := provisioner.NewGoSQL(logger, conf)

	extension := mysql.NewMySQLExtension(prov, conf, logger)

	response, err := extension.GetConnection(request.WorkspaceID, request.ConnectionID)
	if err != nil {
		csmConnection.WriteError(err)
		os.Exit(0)
	}

	csmConnection.Write(*response)
}
