package main

import (
	"os"

	"github.com/hpcloud/sidecar-extensions/go/csm"
	"github.com/hpcloud/sidecar-extensions/go/extension/mongo"
	"github.com/hpcloud/sidecar-extensions/go/extension/mongo/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/mongo/mongoprovisioner"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("mongo-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.MongoDriverConfig{}
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal("main", err)
	}
	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Error("main", err)
		os.Exit(1)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)
	prov := mongoprovisioner.New(conf, logger)

	extension := mongo.NewMongoExtension(prov, conf, logger)

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
