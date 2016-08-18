package db

import (
	"errors"
	"fmt"

	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type SqlDB struct {
	Client Client
}

func NewSqlDB(cfg *config.SqlDB) (*SqlDB, error) {
	if cfg == nil {
		return nil, errors.New("SQL configuration cannot be nil")
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Schema)

	db, err := gorm.Open(cfg.Type, connectionString)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.RouterGroupDB{})
	return &SqlDB{Client: db}, nil
}

func (s *SqlDB) ReadRouterGroups() (models.RouterGroups, error) {
	routerGroupsDB := models.RouterGroupsDB{}
	routerGroups := models.RouterGroups{}
	err := s.Client.Find(&routerGroupsDB).Error
	if err == nil {
		routerGroups = routerGroupsDB.ToRouterGroups()
	}

	return routerGroups, err
}

func (s *SqlDB) ReadRouterGroup(guid string) (models.RouterGroup, error) {
	routerGroupDB := models.RouterGroupDB{}
	routerGroup := models.RouterGroup{}
	err := s.Client.Where("guid = ?", guid).First(&routerGroupDB).Error
	if err == nil {
		routerGroup = routerGroupDB.ToRouterGroup()
	}

	if recordNotFound(err) {
		err = nil
	}

	return routerGroup, err
}

func (s *SqlDB) SaveRouterGroup(routerGroup models.RouterGroup) error {
	existingRouterGroup, err := s.ReadRouterGroup(routerGroup.Guid)
	if err != nil {
		return err
	}

	routerGroupDB := models.NewRouterGroupDB(routerGroup)
	if existingRouterGroup.Guid == routerGroup.Guid {
		updateRouterGroup(&existingRouterGroup, &routerGroup)
		routerGroupDB = models.NewRouterGroupDB(existingRouterGroup)
		err = s.Client.Save(&routerGroupDB).Error
	} else {
		err = s.Client.Create(&routerGroupDB).Error
	}

	return err
}

func updateRouterGroup(existingRouterGroup, currentRouterGroup *models.RouterGroup) {
	if currentRouterGroup.Type != "" {
		existingRouterGroup.Type = currentRouterGroup.Type
	}
	if currentRouterGroup.Name != "" {
		existingRouterGroup.Name = currentRouterGroup.Name
	}
	if currentRouterGroup.ReservablePorts != "" {
		existingRouterGroup.ReservablePorts = currentRouterGroup.ReservablePorts
	}

}
func recordNotFound(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}
