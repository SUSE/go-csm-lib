package main

import (
	"os"

	"github.com/hpcloud/sidecar-extensions/go/csm"
	"github.com/hpcloud/sidecar-extensions/go/extension/redis"
	"github.com/hpcloud/sidecar-extensions/go/extension/redis/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/redis/provisioner"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("redis-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.RedisConfig{}
	env.Parse(&conf)

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Error("main", err)
		os.Exit(1)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath)
	prov := provisioner.NewRedisProvisioner(logger, conf)

	extension := redis.NewRedisExtension(prov, conf, logger)

	response, err := extension.CreateConnection(request.WorkspaceID, request.ConnectionID)
	if err != nil {
		csmConnection.WriteError(err)
		os.Exit(0)
	}

	csmConnection.Write(*response)
}
