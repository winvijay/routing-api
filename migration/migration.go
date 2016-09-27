package migration

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/config"
	"github.com/jinzhu/gorm"
)

const MigrationKey = "routing-api-migration"
const DataVersion = 1

type MigrationData struct {
	MigrationKey   string `gorm:"primary_key"`
	CurrentVersion int
	TargetVersion  int
}

//go:generate counterfeiter -o fakes/fake_migration.go . Migration
type Migration interface {
	RunMigration(db *gorm.DB) error
	Version() int
}

func RunMigrations(sqlCfg *config.SqlDB, migrations []Migration) error {
	// migrations := []Migration{
	// 	NewV0InitMigration(sqlCfg),
	// 	NewV1EtcdMigration(sqlCfg, etcdCfg),
	// }

	if len(migrations) == 0 {
		return nil
	}

	lastMigrationVersion := migrations[len(migrations)-1].Version()

	db, err := connectDB(sqlCfg)
	if err != nil {
		return err
	}
	defer db.Close()
	db.AutoMigrate(&MigrationData{})

	tx := db.Begin()

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
		if existingVersion.TargetVersion < lastMigrationVersion {
			existingVersion.TargetVersion = lastMigrationVersion
			tx.Save(existingVersion)
		}
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	currentVersion := existingVersion.CurrentVersion
	fmt.Printf("\nNumber of migrations: %d\n", len(migrations))
	for _, m := range migrations {
		fmt.Printf("\nLooking at migration version :  %d, CurrentVersion: %d\n", m.Version(), currentVersion)
		if m.Version() > currentVersion {
			fmt.Printf("\nApplying migration version :  %d, Current Version:%d\n", m.Version(), currentVersion)
			m.RunMigration(db)
			currentVersion = m.Version()
			existingVersion.CurrentVersion = currentVersion
			return db.Save(existingVersion).Error
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
