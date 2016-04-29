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
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal("main", err)
		os.Exit(1)
	}

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Fatal("main", err)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)
	prov := provisioner.NewGoSQL(logger, conf)

	extension := mysql.NewMySQLExtension(prov, conf, logger)

	response, err := extension.DeleteWorkspace(request.WorkspaceID)
	if err != nil {
		err := csmConnection.WriteError(err)
		if err != nil {
			logger.Fatal("main", err)
		}
		os.Exit(0)
	}

	err = csmConnection.Write(*response)
	if err != nil {
		logger.Fatal("main", err)
	}
}
