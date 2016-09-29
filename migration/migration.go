package migration

import (
	"fmt"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"github.com/jinzhu/gorm"
)

const MigrationKey = "routing-api-migration"

type MigrationData struct {
	MigrationKey   string `gorm:"primary_key"`
	CurrentVersion int
	TargetVersion  int
}

//go:generate counterfeiter -o fakes/fake_migration.go . Migration
type Migration interface {
	RunMigration(database db.DB) error
	Version() int
}

// TODO Run() {init, RunMigration}
//
func InitializeMigrations(etcdCfg *config.Etcd, logger lager.Logger) []Migration {
	migrations := []Migration{}
	var migration Migration

	migration = NewV0InitMigration()
	migrations = append(migrations, migration)

	// done := make(chan struct{})
	// migration = NewV1EtcdMigration(etcdCfg, done, logger)
	// migrations = append(migrations, migration)

	return migrations
}

func RunMigrations(sqlCfg *config.SqlDB, migrations []Migration) error {
	if len(migrations) == 0 {
		return nil
	}
	fmt.Println("starting migration ************ REMOVE")

	lastMigrationVersion := migrations[len(migrations)-1].Version()

	dbSQL, err := db.NewSqlDB(sqlCfg)
	if err != nil {
		return err
	}

	sqlDB := dbSQL.(*db.SqlDB)
	gormDB := sqlDB.Client.(*gorm.DB)

	defer gormDB.Close()
	gormDB.AutoMigrate(&MigrationData{})

	tx := gormDB.Begin()

	existingVersion := &MigrationData{}

	err = tx.Where("migration_key = ?", MigrationKey).First(existingVersion).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	if err == gorm.ErrRecordNotFound {
		existingVersion = &MigrationData{
			MigrationKey:   MigrationKey,
			CurrentVersion: -1,
			TargetVersion:  lastMigrationVersion,
		}

		err := tx.Create(existingVersion).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if existingVersion.TargetVersion >= lastMigrationVersion {
			return tx.Commit().Error
		}

		existingVersion.TargetVersion = lastMigrationVersion
		err := tx.Save(existingVersion).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	currentVersion := existingVersion.CurrentVersion
	for _, m := range migrations {
		if m.Version() > currentVersion {
			m.RunMigration(dbSQL)
			currentVersion = m.Version()
			existingVersion.CurrentVersion = currentVersion
			err := gormDB.Save(existingVersion).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func connectDB(sqlCfg *config.SqlDB) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		sqlCfg.Username,
		sqlCfg.Password,
		sqlCfg.Host,
		sqlCfg.Port,
		sqlCfg.Schema)

	return gorm.Open(sqlCfg.Type, connectionString)
}
