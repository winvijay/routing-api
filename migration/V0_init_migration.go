package migration

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/jinzhu/gorm"
)

type V0InitMigration struct{}

func NewV0InitMigration() *V0InitMigration {
	return &V0InitMigration{}
}

func (v *V0InitMigration) Version() int {
	return 0
}

func (v *V0InitMigration) RunMigration(dbSQL db.DB) error {
	sqlDB := dbSQL.(*db.SqlDB)
	gormDB := sqlDB.Client.(*gorm.DB)
	fmt.Println("run v0 migration")
	return gormDB.AutoMigrate(&models.RouterGroupDB{}, &models.TcpRouteMapping{}, &models.Route{}).Error
}
