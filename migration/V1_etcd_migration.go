package migration

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/jinzhu/gorm"
)

type V1EtcdMigration struct {
	etcdCfg *config.Etcd
	done    chan struct{}
	logger  lager.Logger
}

func NewV1EtcdMigration(etcdCfg *config.Etcd, done chan struct{}, logger lager.Logger) *V1EtcdMigration {
	return &V1EtcdMigration{etcdCfg: etcdCfg, done: done, logger: logger}
}

func (v *V1EtcdMigration) Version() int {
	return 1
}

func (v *V1EtcdMigration) RunMigration(dbSQL db.DB) error {
	if len(v.etcdCfg.NodeURLS) == 0 {
		v.logger.Info("etcd-not-configured")
		return nil
	}

	sqlDB := dbSQL.(*db.SqlDB)
	gormDB := sqlDB.Client.(*gorm.DB)

	etcd, err := db.NewETCD(*v.etcdCfg)
	if err != nil {
		return err
	}

	etcdRouterGroups, err := etcd.ReadRouterGroups()
	if err != nil {
		return err
	}
	for _, rg := range etcdRouterGroups {
		routerGroup := models.NewRouterGroupDB(rg)
		err := gormDB.Create(&routerGroup).Error
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
		r.ExpiresAt = route.ExpiresAt

		err = gormDB.Create(&r).Error
		if err != nil {
			return err
		}
	}

	fmt.Println("run http events")
	go v.watchForHTTPEvents(etcd, sqlDB)

	etcdTcpRoutes, err := etcd.ReadTcpRouteMappings()
	if err != nil {
		return err
	}
	for _, route := range etcdTcpRoutes {
		r, err := models.NewTcpRouteMappingWithModel(route)
		if err != nil {
			return err
		}
		r.ExpiresAt = route.ExpiresAt

		err = gormDB.Create(&r).Error
		if err != nil {
			return err
		}
	}

	fmt.Println("run tcp events")
	go v.watchForTCPEvents(etcd, sqlDB)

	fmt.Println("run v1 migration")
	return nil
}

func (v *V1EtcdMigration) watchForHTTPEvents(etcd db.DB, sqlDB db.DB) {

	events, errs, cancel := etcd.WatchRouteChanges(db.HTTP_WATCH)
	for {
		select {
		case event := <-events:
			var httpRoute models.Route
			switch event.Type {
			case db.CreateEvent, db.UpdateEvent:
				json.Unmarshal([]byte(event.Value), &httpRoute)
				err := sqlDB.SaveRoute(httpRoute)
				if err != nil {
					v.logger.Error("failed-to-save-http-route", err)
				}
			case db.DeleteEvent:
				json.Unmarshal([]byte(event.Value), &httpRoute)
				err := sqlDB.DeleteRoute(httpRoute)
				if err != nil {
					v.logger.Error("failed-to-delete-http-route", err)
				}
			default:
				break
			}
		case <-errs:
			break
		case <-v.done:
			cancel()
			return
		}
	}
}

func (v *V1EtcdMigration) watchForTCPEvents(etcd db.DB, sqlDB db.DB) {

	events, errs, cancel := etcd.WatchRouteChanges(db.TCP_WATCH)
	for {
		select {
		case event := <-events:
			var tcpRoute models.TcpRouteMapping
			switch event.Type {
			case db.CreateEvent, db.UpdateEvent:
				json.Unmarshal([]byte(event.Value), &tcpRoute)
				err := sqlDB.SaveTcpRouteMapping(tcpRoute)
				if err != nil {
					v.logger.Error("failed-to-save-tcp-route", err)
				}
			case db.DeleteEvent:
				json.Unmarshal([]byte(event.Value), &tcpRoute)
				err := sqlDB.DeleteTcpRouteMapping(tcpRoute)
				if err != nil {
					v.logger.Error("failed-to-delete-tcp-route", err)
				}
			default:
				break
			}
		case <-errs:
			break
		case <-v.done:
			cancel()
			return
		}
	}
}
