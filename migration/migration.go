package migration

import "code.cloudfoundry.org/routing-api/config"

type Migration interface {
	RunMigration() error
	Version() int
}

func RunMigrations(sqlCfg config.SqlDB, etcdCfg config.Etcd) {
	migrations := []Migration{
		NewV0InitMigration(sqlCfg),
		new(V1EtcdMigration),
	}

	for _, m := range migrations {
		m.RunMigration()
	}
}
