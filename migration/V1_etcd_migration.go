package migration

import "code.cloudfoundry.org/routing-api/db"

// one-time migration

type V1EtcdMigration struct {
	sqlDB db.DB
	etcd  db.DB
}

func NewV1EtcdMigration(sqlDB db.DB, etcd db.DB) *V1EtcdMigration {
	return &V1EtcdMigration{sqlDB: sqlDB, etcd: etcd}
}

func (v *V1EtcdMigration) Version() int {
	return 1
}

func (v *V1EtcdMigration) RunMigration() error {
	// connect to etcd
	etcdRouterGroups, err := v.etcd.ReadRouterGroups()
	if err != nil {
		return err
	}
	for _, rg := range etcdRouterGroups {
		err := v.sqlDB.SaveRouterGroup(rg)
		if err != nil {
			return err
		}
	}
	return nil
}
