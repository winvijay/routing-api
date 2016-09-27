package migration

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/jinzhu/gorm"
)

type V0InitMigration struct {
	sqlCfg *config.SqlDB
}

func NewV0InitMigration(cfg *config.SqlDB) *V0InitMigration {
	return &V0InitMigration{sqlCfg: cfg}
}

func (v *V0InitMigration) Version() int {
	return 0
}

func (v *V0InitMigration) RunMigration() error {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		v.sqlCfg.Username,
		v.sqlCfg.Password,
		v.sqlCfg.Host,
		v.sqlCfg.Port,
		v.sqlCfg.Schema)

	db, err := gorm.Open(v.sqlCfg.Type, connectionString)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.AutoMigrate(&models.RouterGroupDB{}, &models.TcpRouteMapping{}, &models.Route{}).Error
}
