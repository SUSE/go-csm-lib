package main

import (
	"os"

	"github.com/hpcloud/sidecar-extensions/go/csm"
	"github.com/hpcloud/sidecar-extensions/go/extension/rabbitmq"
	"github.com/hpcloud/sidecar-extensions/go/extension/rabbitmq/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/rabbitmq/provisioner"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("rabbitmq-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.RabbitmqConfig{}
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal("main", err)
	}
	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Fatal("main", err)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)
	prov := provisioner.NewRabbitHoleProvisioner(logger, conf)

	extension := rabbitmq.NewRabbitmqExtension(prov, conf, logger)

	response, err := extension.GetWorkspace(request.WorkspaceID)
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
