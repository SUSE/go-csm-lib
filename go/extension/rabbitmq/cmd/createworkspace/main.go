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
	env.Parse(&conf)

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Error("main", err)
		os.Exit(1)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath)
	prov := provisioner.NewRabbitHoleProvisioner(logger, conf)

	extension := rabbitmq.NewRabbitmqExtension(prov, conf, logger)

	response, err := extension.CreateWorkspace(request.WorkspaceID)
	if err != nil {
		csmConnection.WriteError(err)
		os.Exit(0)
	}

	csmConnection.Write(*response)
}
