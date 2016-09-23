package migration

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/jinzhu/gorm"
)

type V0InitMigration struct {
	sqlDB *gorm.DB
}

func NewV0InitMigration(cfg config.SqlDB) *V0InitMigration {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Schema)

	db, err := gorm.Open(cfg.Type, connectionString)
	if err != nil {
		return &V0InitMigration{}
	}

	return &V0InitMigration{sqlDB: db}
}

func (v *V0InitMigration) Version() int {
	return 0
}

func (v *V0InitMigration) RunMigration() error {
	return v.sqlDB.AutoMigrate(&models.RouterGroupDB{}, &models.TcpRouteMapping{}, &models.Route{}).Error
}
