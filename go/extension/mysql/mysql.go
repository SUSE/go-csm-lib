package mysql

import (
	"github.com/hpcloud/sidecar-extensions/go/extension"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/provisioner"
)

type mysqlExtension struct {
	conf        config.MySQLConfig
	provisioner provisioner.MySQLProvisioner
}

func NewMySQLExtension(provisioner provisioner.MySQLProvisioner, conf config.MySQLConfig) extension.Extension {
	return &mysqlExtension{conf: conf, provisioner: provisioner}
}
