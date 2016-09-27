package migration

import (
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"
)

type V1EtcdMigration struct {
	sqlCfg  *config.SqlDB
	etcdCfg *config.Etcd
}

func NewV1EtcdMigration(sqlCfg *config.SqlDB, etcdCfg *config.Etcd) *V1EtcdMigration {
	return &V1EtcdMigration{sqlCfg: sqlCfg, etcdCfg: etcdCfg}
}

func (v *V1EtcdMigration) Version() int {
	return 1
}

func (v *V1EtcdMigration) RunMigration() error {
	s, err := db.NewSqlDB(v.sqlCfg)
	if err != nil {
		return err
	}
	sqlDB := s.(*db.SqlDB)

	etcd, err := db.NewETCD(*v.etcdCfg)
	if err != nil {
		return err
	}

	etcdRouterGroups, err := etcd.ReadRouterGroups()
	if err != nil {
		return err
	}
	for _, rg := range etcdRouterGroups {
		err := sqlDB.SaveRouterGroup(rg)
		if err != nil {
			return err
		}
	}

	etcdRoutes, err := etcd.ReadRoutes()
	if err != nil {
		return err
	}
	for _, route := range etcdRoutes {
		r, err := models.NewRouteWithModel(route)
		if err != nil {
			return err
		}

		err = sqlDB.Client.Create(&r).Error
		if err != nil {
			return err
		}
	}

	etcdTcpRoutes, err := etcd.ReadTcpRouteMappings()
	if err != nil {
		return err
	}
	for _, route := range etcdTcpRoutes {
		r, err := models.NewTcpRouteMappingWithModel(route)
		if err != nil {
			return err
		}

		err = sqlDB.Client.Create(&r).Error
		if err != nil {
			return err
		}
	}
	return nil
}
